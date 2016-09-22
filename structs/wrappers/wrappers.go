package wrappers

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"time"
)

// log wrapper logs every time a request is made
type LogWrapper struct {
	client.Client
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

// Implements client.Wrapper as logWrapper
func LogWrap(c client.Client) client.Client {
	return &LogWrapper{c}
}

func NewKazoupClient() client.Client {
	c := client.NewClient(
		client.Wrap(LogWrap),
	)
	return c
}

func NewKazoupService(name string) micro.Service {
	sn := fmt.Sprintf("%s.srv.%s", globals.NAMESPACE, name)

	service := micro.NewService(
		micro.Name(sn),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
		micro.Client(NewKazoupClient()),
	)
	return service
}
