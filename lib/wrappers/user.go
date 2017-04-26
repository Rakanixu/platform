package wrappers

import (
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

const (
	TEST_USER_ID = "test_user"
)

// UserTestHandlerWrapper sets a user in the context for tests pusposes
func UserTestHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		ctx = context.WithValue(
			ctx,
			kazoup_context.UserIdCtxKey{},
			kazoup_context.UserIdCtxValue(TEST_USER_ID),
		)

		return fn(ctx, req, rsp)
	}
}

// UserTestSubscriberWrapper sets a user in the context for tests pusposes
func UserTestSubscriberWrapper(fn server.SubscriberFunc) server.SubscriberFunc {
	return func(ctx context.Context, msg server.Publication) error {
		ctx = context.WithValue(
			ctx,
			kazoup_context.UserIdCtxKey{},
			kazoup_context.UserIdCtxValue(TEST_USER_ID),
		)

		return fn(ctx, msg)
	}
}
