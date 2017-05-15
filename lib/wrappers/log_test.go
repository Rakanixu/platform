package wrappers

import (
	"errors"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"testing"
)

func TestNewLogHandlerWrapper(t *testing.T) {
	NewLogHandlerWrapper()
}

func Test_logHandlerWrapper(t *testing.T) {
	callbackError := errors.New("Log error")

	fn := func(h server.HandlerFunc) server.HandlerFunc {
		return logHandlerWrapper(h)
	}

	handlerFuncSuccess := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return nil
	}

	handlerFuncFail := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return callbackError
	}

	testData := []struct {
		ctx         context.Context
		handlerFunc server.HandlerFunc
		req         server.Request
		rsp         interface{}
		err         error
	}{
		{
			ctx:         context.TODO(),
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         nil,
		},
		{
			ctx:         context.TODO(),
			handlerFunc: handlerFuncFail,
			req:         request{},
			rsp:         response{},
			err:         callbackError,
		},
	}

	for _, tt := range testData {
		result := fn(tt.handlerFunc)(tt.ctx, tt.req, tt.rsp)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}

func TestNewLogSubscriberWrapper(t *testing.T) {
	NewLogSubscriberWrapper()
}

func Test_logSubscriberWrapper(t *testing.T) {
	callbackError := errors.New("Log error")

	fn := func(h server.SubscriberFunc) server.SubscriberFunc {
		return logSubscriberWrapper(h)
	}

	handlerFuncSuccess := func(ctx context.Context, msg server.Publication) error {
		return nil
	}

	handlerFuncFail := func(ctx context.Context, msg server.Publication) error {
		return callbackError
	}

	testData := []struct {
		ctx         context.Context
		handlerFunc server.SubscriberFunc
		msg         server.Publication
		err         error
	}{
		{
			ctx:         ctx,
			handlerFunc: handlerFuncSuccess,
			msg:         &publication{},
			err:         nil,
		},
		{
			ctx:         context.TODO(),
			handlerFunc: handlerFuncFail,
			msg:         &publication{},
			err:         callbackError,
		},
	}

	for _, tt := range testData {
		result := fn(tt.handlerFunc)(tt.ctx, tt.msg)

		if tt.err != result {
			t.Errorf("Expected %v, got %v", tt.err, result)
		}
	}
}
