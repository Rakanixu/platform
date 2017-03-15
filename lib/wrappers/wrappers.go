package wrappers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/dgrijalva/jwt-go"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"github.com/micro/go-plugins/wrapper/trace/awsxray"
	"golang.org/x/net/context"
	timerate "golang.org/x/time/rate"
	"gopkg.in/go-redis/rate.v5"
	"gopkg.in/redis.v5"
	"os"
	"time"
)

type kazoupClientWrapper struct {
	client.Client
}

// KazoupClientWrap wraps client
func KazoupClientWrap() client.Wrapper {
	return func(c client.Client) client.Client {
		return &kazoupClientWrapper{c}
	}
}

// Call will set X-Kazoup-Token with DB_ACCESS_TOKEN value in every internal request
// After every call, we will publish an announcment to say what happened
func (kcw *kazoupClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	md["X-Kazoup-Token"] = globals.DB_ACCESS_TOKEN
	ctx = metadata.NewContext(ctx, md)

	// Execute the call
	if err := kcw.Client.Call(ctx, req, rsp, opts...); err != nil {
		return err
	}

	b, err := json.Marshal(req.Request())
	if err != nil {
		return err
	}

	// Publish annuncment after handler was called
	if err := kcw.Client.Publish(ctx, kcw.Client.NewPublication(
		globals.AnnounceTopic,
		&announce.AnnounceMessage{
			Handler: fmt.Sprintf("%s.%s", req.Service(), req.Method()),
			Data:    string(b),
		},
	)); err != nil {
		return err
	}

	return nil
}

// NewHandlerWrapper wraps a service within the handler so it can be accessed by the handler itself.
func NewHandlerWrapper(service micro.Service) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = micro.NewContext(ctx, service)
			return h(ctx, req, rsp)
		}
	}
}

// log wrapper logs every time a request is made
type LogWrapper struct {
	client.Client
}

// Implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &LogWrapper{c}
}

func (l *LogWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	log.WithFields(log.Fields{
		"from":    md["X-Micro-From-Service"],
		"service": req.Service(),
		"method":  req.Method(),
		"request": req.Request(),
	}).Info("Service call")
	return l.Client.Call(ctx, req, rsp, opts...)
}

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var f error

		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.InternalServerError("AuthWrapper", "Unable to retrieve metadata")
		}

		if len(md["Authorization"]) == 0 {
			return errors.Unauthorized("", "Authorization required")
		}

		// Authentication
		if md["Authorization"] != globals.SYSTEM_TOKEN {
			token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
				if err != nil {
					return nil, err
				}

				return decoded, nil
			})

			if err != nil {
				return errors.Unauthorized("Token", err.Error())
			}

			if !token.Valid {
				return errors.Unauthorized("", "Invalid token")
			}

			ctx = context.WithValue(
				ctx,
				kazoup_context.UserIdCtxKey{},
				kazoup_context.UserIdCtxValue(token.Claims.(jwt.MapClaims)["sub"].(string)),
			)
		}

		f = fn(ctx, req, rsp)

		return f
	}
}

// quotaHandlerWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func quotaHandlerWrapper(fn server.HandlerFunc, limiter *rate.Limiter, srv string) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		var f error

		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.InternalServerError("QuotaWrapper", "Unable to retrieve metadata")
		}

		if len(md["Authorization"]) == 0 {
			return errors.Unauthorized("QuotaWrapper", "No Auth header")
		}

		// We will read claim to know if public user, or paying or whatever
		token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
			if err != nil {
				return nil, err
			}

			return decoded, nil
		})
		if err != nil {
			return errors.Unauthorized("Token", err.Error())
		}

		if token.Claims.(jwt.MapClaims)["roles"] == nil {
			return errors.BadRequest("Roles not found", "Unable to retrieve user roles.")
		}

		var quotaLimit int64
		for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
			switch v.(string) {
			case globals.PRODUCT_TYPE_PERSONAL, globals.PRODUCT_TYPE_TEAM, globals.PRODUCT_TYPE_ENTERPRISE:
				quotaLimit = int64(globals.PRODUCT_QUOTAS.M[v.(string)][srv]["handler"].(int))
			}
		}

		if quotaLimit <= 0 {
			return fn(ctx, req, rsp)
		}

		_, _, allowed := limiter.AllowN(fmt.Sprintf("%s-handler-%s", srv, token.Claims.(jwt.MapClaims)["sub"].(string)), quotaLimit, globals.QUOTA_TIME_LIMITER, 1)
		if !allowed {
			return errors.Forbidden("User Rate Limit", "User rate limit exceeded.")
		}

		f = fn(ctx, req, rsp)

		return f
	}
}

// NewQuotaHandlerWrapper returns a handler quota limit per user wrapper
func NewQuotaHandlerWrapper(srvName string) server.HandlerWrapper {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "redis:6379",
		},
	})
	fallbackLimiter := timerate.NewLimiter(timerate.Every(time.Second), 1000)
	limiter := rate.NewLimiter(ring, fallbackLimiter)

	return func(h server.HandlerFunc) server.HandlerFunc {
		return quotaHandlerWrapper(h, limiter, srvName)
	}
}

// quotaSubscriberWrapper defines a quota wrapper based on quotaLimit per srv+user_id key
func quotaSubscriberWrapper(fn server.SubscriberFunc, limiter *rate.Limiter, srv string) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.InternalServerError("QuotaWrapper", "Unable to retrieve metadata")
		}

		if len(md["Authorization"]) == 0 {
			return errors.Unauthorized("QuotaWrapper", "No Auth header")
		}

		// We will read claim to know if public user, or paying or whatever
		token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			decoded, err := base64.URLEncoding.DecodeString(globals.CLIENT_ID_SECRET)
			if err != nil {
				return nil, err
			}

			return decoded, nil
		})
		if err != nil {
			return errors.Unauthorized("Token", err.Error())
		}

		if token.Claims.(jwt.MapClaims)["roles"] == nil {
			return errors.BadRequest("Roles not found", "Unable to retrieve user roles.")
		}

		var quotaLimit int64
		for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
			switch v.(string) {
			case globals.PRODUCT_TYPE_PERSONAL, globals.PRODUCT_TYPE_TEAM, globals.PRODUCT_TYPE_ENTERPRISE:
				quotaLimit = int64(globals.PRODUCT_QUOTAS.M[v.(string)][srv]["subscriber"].(int))
			}
		}

		if quotaLimit <= 0 {
			return fn(ctx, msg)
		}

		_, _, allowed := limiter.AllowN(fmt.Sprintf("%s-subs-%s", srv, token.Claims.(jwt.MapClaims)["sub"].(string)), quotaLimit, globals.QUOTA_TIME_LIMITER, 1)
		if !allowed {
			// Quota limite reached, but due to subscribers nature, error will be lost.
			// IDEA: pulbish to notification srv a rate limite message to let user know.
			log.Println("USER RATE LIMIT (SUBSCRIBER)", fmt.Sprintf("%s%s", srv, token.Claims.(jwt.MapClaims)["sub"].(string)))
			return errors.Forbidden("User Rate Limit", "User rate limit exceeded.")
		}

		return fn(ctx, msg)
	}
}

// NewQuotaSubscriberWrapper returns a subscriber quota limit per user wrapper
func NewQuotaSubscriberWrapper(srvName string) server.SubscriberWrapper {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": "redis:6379",
		},
	})
	fallbackLimiter := timerate.NewLimiter(timerate.Every(time.Second), 1000)
	limiter := rate.NewLimiter(ring, fallbackLimiter)

	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return quotaSubscriberWrapper(fn, limiter, srvName)
	}
}

func NewKazoupClient() client.Client {
	return client.NewClient()
}

func NewKazoupClientWithXrayTrace(sess *session.Session) client.Client {
	opts := []awsxray.Option{
		// Used as segment name
		awsxray.WithName("com.kazoup.client"),
		// Specify X-Ray Daemon Address
		// awsxray.WithDaemon("localhost:2000"),
		// Or X-Ray Client
		awsxray.WithClient(xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})),
	}
	return client.NewClient(
		client.WrapCall(awsxray.NewCallWrapper(opts...)),
	)
}

func NewKazoupService(name string, mntr ...monitor.Monitor) micro.Service {
	var m monitor.Monitor

	// Check if monitor available
	if len(mntr) > 0 && mntr[0] != nil {
		m = mntr[0]
	}

	//Get AWS session credentials in ~/.aws/credentials
	/*	sess, err := session.NewSession()
		if err != nil {
			log.Fatal(err)
		}
		// New AWSXRAY with default
		//x := xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
		opts := []awsxray.Option{
			// Used as segment name
			awsxray.WithName(name),
			// Specify X-Ray Daemon Address
			// awsxray.WithDaemon("localhost:2000"),
			// Or X-Ray Client
			awsxray.WithClient(xray.New(sess, &aws.Config{Region: aws.String("eu-west-1")})),
		}*/

	//FIXME hacked we should just pass micro.Flags
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("oops can't get hostname")
	}

	md := map[string]string{
		"hostname": hostname,
	}

	sn := fmt.Sprintf("%s.srv.%s", globals.NAMESPACE, name)

	if name == "db" {
		service := micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
			//micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapClient( /*awsxray.NewClientWrapper(opts...),*/ KazoupClientWrap()),
			micro.WrapSubscriber(NewQuotaSubscriberWrapper(sn)),
			micro.WrapHandler( /*awsxray.NewHandlerWrapper(opts...), */ NewQuotaHandlerWrapper(sn), AuthWrapper),
			micro.Flags(
				cli.StringFlag{
					Name:   "elasticsearch_hosts",
					EnvVar: "ELASTICSEARCH_HOSTS",
					Usage:  "Comma separated list of elasticsearch hosts",
					Value:  "elasticsearch:9200",
				},
			),
			micro.Action(func(c *cli.Context) {
				//parts := strings.Split(c.String("elasticsearch_hosts"), ",")
				//elastic.Hosts = parts
			}),
		)
		return service
	}

	var service micro.Service
	if m == nil {
		service = micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
			//micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapClient( /*awsxray.NewClientWrapper(opts...),*/ KazoupClientWrap()),
			micro.WrapSubscriber(NewQuotaSubscriberWrapper(sn)),
			micro.WrapHandler( /*awsxray.NewHandlerWrapper(opts...), */ NewQuotaHandlerWrapper(sn), AuthWrapper),
		)
	} else {
		service = micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
			//micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapClient( /*awsxray.NewClientWrapper(opts...), */ monitor.ClientWrapper(m), KazoupClientWrap()),
			micro.WrapSubscriber(NewQuotaSubscriberWrapper(sn)),
			micro.WrapHandler( /*awsxray.NewHandlerWrapper(opts...),*/ monitor.HandlerWrapper(m), NewQuotaHandlerWrapper(sn), AuthWrapper),
		)
	}

	return service
}
