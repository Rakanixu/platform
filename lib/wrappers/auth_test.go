package wrappers

import (
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"testing"
)

func TestNewAuthHandlerWrapper(t *testing.T) {
	NewAuthHandlerWrapper()
}

func Test_authHandlerWrapper(t *testing.T) {
	fn := func(h server.HandlerFunc) server.HandlerFunc {
		return authHandlerWrapper(h)
	}

	handlerFuncSuccess := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return nil
	}

	// Succesfull authorization requires not expired token
	testData := []struct {
		ctx         context.Context
		handlerFunc server.HandlerFunc
		req         server.Request
		rsp         interface{}
		err         error
	}{
		// No metadata
		{
			ctx:         micro.NewContext(ctx, srv),
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         platform_errors.ErrInvalidCtx,
		},
		// Empty authorization header
		{
			ctx: micro.NewContext(metadata.NewContext(ctx, map[string]string{
				"Authorization": "",
			}), srv),
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         platform_errors.ErrNoAuthHeader,
		},
		/*		{
				ctx: micro.NewContext(metadata.NewContext(ctx, map[string]string{
					"Authorization": "123",
				}), srv),
				handlerFunc: handlerFuncSuccess,
				req:         request{},
				rsp:         response{},
				err:         platform_errors.NewPlatformError("", "ParseJWTToken", "", errors.New("token contains an invalid number of segments")),
			},*/
	}

	for _, tt := range testData {
		result := fn(tt.handlerFunc)(tt.ctx, tt.req, tt.rsp)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}

func Test_authSubscriberWrapper(t *testing.T) {
	fn := func(h server.SubscriberFunc) server.SubscriberFunc {
		return authSubscriberWrapper(h)
	}

	handlerSubscriberSuccess := func(ctx context.Context, msg server.Publication) error {
		return nil
	}

	testData := []struct {
		ctx            context.Context
		subscriberFunc server.SubscriberFunc
		msg            server.Publication
		err            error
	}{
		{
			ctx:            micro.NewContext(ctx, srv),
			subscriberFunc: handlerSubscriberSuccess,
			msg:            publication{},
			err:            platform_errors.ErrInvalidCtx,
		},
		{
			ctx: micro.NewContext(metadata.NewContext(ctx, map[string]string{
				"Authorization": "",
			}), srv),
			subscriberFunc: handlerSubscriberSuccess,
			msg:            publication{},
			err:            platform_errors.ErrNoAuthHeader,
		},
	}

	for _, tt := range testData {
		result := fn(tt.subscriberFunc)(tt.ctx, tt.msg)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}
