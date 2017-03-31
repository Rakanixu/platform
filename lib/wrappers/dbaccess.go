package wrappers

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
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

	return kcw.Client.Call(ctx, req, rsp, opts...)
}
