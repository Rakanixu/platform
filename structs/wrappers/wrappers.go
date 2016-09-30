package wrappers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kazoup/platform/structs/globals"
	auth_proto "github.com/micro/auth-srv/proto/oauth2"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"os"
	"os/user"
	"time"
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

		if len(md["Token"]) == 0 {
			return errors.Unauthorized("", "")
		}

		if md["Token"] != globals.SYSTEM_TOKEN {
			c := auth_proto.NewOauth2Client(globals.AUTH_SERVICE_NAME, nil)
			r, err := c.Introspect(ctx, &auth_proto.IntrospectRequest{
				AccessToken: md["Token"],
			})
			if err != nil {
				return errors.InternalServerError("AuthWrapper", err.Error())
			}

			if r.GetToken() == nil {
				return errors.Unauthorized("", "")
			}
		}

		/*		if !r.Active {
					cfg := &oauth2.Config{}
					tokenSource := cfg.TokenSource(oauth2.NoContext, &oauth2.Token{
						AccessToken:  r.Token.AccessToken,
						TokenType:    r.Token.TokenType,
						RefreshToken: r.Token.RefreshToken,
						Expiry:       time.Unix(r.Token.ExpiresAt, 0),
					})
					t, err := tokenSource.Token()
					if err != nil {
						return errors.InternalServerError("AuthWrapper", err.Error())
					}

					//newCtx := client.NewContext(ctx, t)

					newCtx := metadata.NewContext(ctx, map[string]string{
						"Token": t.AccessToken,
					})

					log.Println("========")
					log.Println(newCtx)

					f = fn(newCtx, req, rsp)
				} else {

				}*/
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
	c := client.NewClient(
		client.Wrap(DesktopWrap),
	)
	return c
}

func NewKazoupService(name string) micro.Service {
	//FIXME hacked we should just pass micro.Flags
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("oops can't get hostname")
	}
	u, err := user.Current()
	if err != nil {
		log.Fatal("oops can't get username")
	}

	md := map[string]string{
		"hostname": hostname,
		"username": u.Username,
	}

	if name == "db" {
		sn := fmt.Sprintf("%s.srv.%s", globals.NAMESPACE, name)
		service := micro.NewService(
			micro.Name(sn),
			micro.Version("latest"),
			micro.Metadata(md),
			micro.Client(NewKazoupClient()),
			micro.WrapSubscriber(SubscriberWrapper),
			micro.WrapHandler(AuthWrapper),
			micro.Flags(
				cli.StringFlag{
					Name:   "elasticsearch_hosts",
					EnvVar: "ELASTICSEARCH_HOSTS",
					Usage:  "Comma separated list of elasticsearch hosts",
					Value:  "localhost:9200",
				},
			),
			micro.Action(func(c *cli.Context) {
				//parts := strings.Split(c.String("elasticsearch_hosts"), ",")
				//elastic.Hosts = parts
			}),
		)

		return service
	}
	sn := fmt.Sprintf("%s.srv.%s", globals.NAMESPACE, name)
	service := micro.NewService(
		micro.Name(sn),
		micro.Version("latest"),
		micro.Metadata(md),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
		micro.Client(NewKazoupClient()),
		micro.WrapSubscriber(SubscriberWrapper),
		micro.WrapHandler(AuthWrapper),
	)
	return service
}
