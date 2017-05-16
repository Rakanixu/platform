package wrappers

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

type contextClientWrapper struct {
	service micro.Service
	client.Client
}

func (c *contextClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	ctx = micro.NewContext(ctx, c.service)

	return c.Client.Call(ctx, req, rsp, opts...)
}

func ContextClientWrapper(service micro.Service) client.Wrapper {
	return func(c client.Client) client.Client {
		return &contextClientWrapper{service, c}
	}
}

func NewContextHandlerWrapper(service micro.Service) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return contextHandlerWrapper(h, service)
	}
}

func contextHandlerWrapper(h server.HandlerFunc, srv micro.Service) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		ctx = micro.NewContext(ctx, srv)

		return h(ctx, req, rsp)
	}
}

func NewContextSubscriberWrapper(service micro.Service) server.SubscriberWrapper {
	return func(fn server.SubscriberFunc) server.SubscriberFunc {
		return contextSubscriberWrapper(fn, service)
	}
}

func contextSubscriberWrapper(fn server.SubscriberFunc, srv micro.Service) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		ctx = micro.NewContext(ctx, srv)

		return fn(ctx, msg)
	}
}
