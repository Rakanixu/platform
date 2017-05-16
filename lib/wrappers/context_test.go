package wrappers

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"reflect"
	"testing"
)

func TestContextClientWrapper_Call(t *testing.T) {
	clientWrapper := &contextClientWrapper{
		srv,
		client.NewClient(),
	}

	clientWrapper.Call(context.TODO(), request{}, response{})
}

func TestContextClientWrapper(t *testing.T) {
	ContextClientWrapper(srv)
}

func TestNewContextHandlerWrapper(t *testing.T) {
	NewContextHandlerWrapper(srv)
}

func Test_contextHandlerWrapper(t *testing.T) {
	fn := func(h server.HandlerFunc) server.HandlerFunc {
		return contextHandlerWrapper(h, srv)
	}

	handlerFuncSuccess := func(ctx context.Context, req server.Request, rsp interface{}) error {
		service, ok := micro.FromContext(ctx)
		if !ok {
			t.Fatal()
		}

		if !reflect.DeepEqual(srv, service) {
			t.Error("Service in context not equal to service passed")
		}

		return nil
	}

	fn(handlerFuncSuccess)(context.TODO(), request{}, response{})
}

func TestNewContextSubscriberWrapper(t *testing.T) {
	NewContextSubscriberWrapper(srv)
}

func Test_contextSubscriberWrapper(t *testing.T) {
	fn := func(h server.SubscriberFunc) server.SubscriberFunc {
		return contextSubscriberWrapper(h, srv)
	}

	subscriberFuncSuccess := func(ctx context.Context, msg server.Publication) error {
		service, ok := micro.FromContext(ctx)
		if !ok {
			t.Fatal()
		}

		if !reflect.DeepEqual(srv, service) {
			t.Error("Service in context not equal to service passed")
		}

		return nil
	}

	fn(subscriberFuncSuccess)(context.TODO(), publication{})
}
