package wrappers

import (
	kazoup_context "github.com/kazoup/platform/lib/context"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"testing"
)

var (
	srv = NewKazoupService("test-service")
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

/* Helpers to mock micro handlers and subscribers */

type request struct{}

func (r request) Service() string {
	return "service"
}

func (r request) Method() string {
	return "method"
}

func (r request) ContentType() string {
	return "Content-Type"
}

func (r request) Request() interface{} {
	return r
}

func (r request) Stream() bool {
	return false
}

type response struct{}

type publication struct{}

func (p publication) Topic() string {
	return "topic"
}

func (p publication) Message() interface{} {
	type msg struct{}

	return new(msg)
}

func (p publication) ContentType() string {
	return "Content-Type"
}

func TestNewAfterHandlerWrapper(t *testing.T) {
	NewAfterHandlerWrapper()
}

func Test_afterHandlerWrapper(t *testing.T) {
	handlerError := errors.New("handlerError")

	fn := func(h server.HandlerFunc) server.HandlerFunc {
		return afterHandlerWrapper(h)
	}

	handlerFuncSuccess := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return nil
	}

	handlerFuncFail := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return handlerError
	}

	testData := []struct {
		ctx         context.Context
		handlerFunc server.HandlerFunc
		req         server.Request
		rsp         interface{}
		err         error
	}{
		{
			ctx:         micro.NewContext(ctx, srv),
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         nil,
		},
		{
			ctx:         micro.NewContext(ctx, srv),
			handlerFunc: handlerFuncFail,
			req:         request{},
			rsp:         response{},
			err:         handlerError,
		},
		{
			ctx:         ctx,
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range testData {
		result := fn(tt.handlerFunc)(tt.ctx, tt.req, tt.rsp)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}

func TestNewAfterSubscriberWrapper(t *testing.T) {
	NewAfterSubscriberWrapper()
}

func Test_afterSubscriberWrapper(t *testing.T) {
	subscriberError := errors.New("subscriberError")

	fn := func(h server.SubscriberFunc) server.SubscriberFunc {
		return afterSubscriberWrapper(h)
	}

	handlerSubscriberSuccess := func(ctx context.Context, msg server.Publication) error {
		return nil
	}

	handlerSubscriberFail := func(ctx context.Context, msg server.Publication) error {
		return subscriberError
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
			err:            nil,
		},
		{
			ctx:            micro.NewContext(ctx, srv),
			subscriberFunc: handlerSubscriberFail,
			msg:            publication{},
			err:            subscriberError,
		},
		{
			ctx:            ctx,
			subscriberFunc: handlerSubscriberSuccess,
			msg:            publication{},
			err:            platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range testData {
		result := fn(tt.subscriberFunc)(tt.ctx, tt.msg)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}
