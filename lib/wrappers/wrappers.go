package wrappers

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/xray"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"github.com/micro/go-plugins/wrapper/trace/awsxray"
	"golang.org/x/net/context"
)

// log wrapper logs every time a request is made
type LogWrapper struct {
	client.Client
}

type DesktopWrapper struct {
	client.Client
}

type clientWrapper struct {
	service micro.Service
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = micro.NewContext(ctx, c.service)
	return c.Client.Call(ctx, req, rsp, opts...)
}

// NewClientWrapper wraps a service within a client so it can be accessed by subsequent client wrappers.
func NewClientWrapper(service micro.Service) client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{service, c}
	}
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

func (dw *DesktopWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	hostname, _ := os.Hostname()
	md, _ := metadata.FromContext(ctx)
	filter := func(services []*registry.Service) []*registry.Service {
		for _, service := range services {
			var nodes []*registry.Node
			for _, node := range service.Nodes {
				if node.Metadata["hostname"] == hostname {
					nodes = append(nodes, node)
				}
			}
			service.Nodes = nodes
		}
		return services
	}
	callOptions := append(opts, client.WithSelectOption(
		selector.WithFilter(filter),
	))
	log.WithFields(log.Fields{
		"hostname": hostname,
		"from":     md["X-Micro-From-Service"],
		"service":  req.Service(),
		"method":   req.Method(),
		"request":  req.Request(),
	}).Info("[Dekstop Wrapper] filtering for hostname")
	return dw.Client.Call(ctx, req, rsp, callOptions...)
}

func (l *LogWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	md, _ := metadata.FromContext(ctx)
	log.WithFields(log.Fields{
		"from":    md["X-Micro-From-Service"],
		"service": req.Service(),
		"method":  req.Method(),
		"request": req.Request(),
	}).Info("Service call")
	return l.Client.Call(ctx, req, rsp)
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

			// Authorization, inject id in context
			ctx = metadata.NewContext(context.TODO(), map[string]string{
				"Authorization": md["Authorization"],
				"Id":            token.Claims.(jwt.MapClaims)["sub"].(string),
			})
		}

		f = fn(ctx, req, rsp)

		return f
	}
}

// SubscriberWrapper for auth internal async messages
func SubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		var f error

		f = fn(globals.NewSystemContext(), msg)

		return f
	}
}

// Implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &LogWrapper{c}
}
func DesktopWrap(c client.Client) client.Client {
	return &DesktopWrapper{c}
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
	sess, err := session.NewSession()
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
	}

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
			micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapClient(awsxray.NewClientWrapper(opts...)),
			micro.WrapSubscriber(SubscriberWrapper),
			micro.WrapHandler(awsxray.NewHandlerWrapper(opts...), AuthWrapper),
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
			micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapSubscriber(SubscriberWrapper),
			micro.WrapClient(awsxray.NewClientWrapper(opts...)),
			micro.WrapHandler(awsxray.NewHandlerWrapper(opts...), AuthWrapper),
		)
	} else {
		service = micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.RegisterTTL(time.Minute),
			micro.RegisterInterval(time.Second*30),
			micro.Client(NewKazoupClientWithXrayTrace(sess)),
			micro.WrapSubscriber(SubscriberWrapper),
			micro.WrapClient(awsxray.NewClientWrapper(opts...), monitor.ClientWrapper(m)),
			micro.WrapHandler(awsxray.NewHandlerWrapper(opts...), monitor.HandlerWrapper(m), AuthWrapper),
		)
	}

	return service
}
