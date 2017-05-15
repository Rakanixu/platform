package wrappers

import (
	/*	"github.com/go-redis/redis"
		rate "github.com/go-redis/redis_rate"
		"github.com/kazoup/platform/lib/errors"
		"github.com/micro/go-micro"
		"github.com/micro/go-micro/metadata"
		"github.com/micro/go-micro/server"
		"golang.org/x/net/context"
		timerate "golang.org/x/time/rate"*/
	"testing"
	/*	"time"*/)

func TestNewQuotaHandlerWrapper(t *testing.T) {
	NewQuotaHandlerWrapper("service-name")
}

// Tests requires a valid JWT token, but it will expired
/*
func Test_quotaHandlerWrapper(t *testing.T) {
	fn := func(h server.HandlerFunc) server.HandlerFunc {
		ring := redis.NewRing(&redis.RingOptions{
			Addrs: map[string]string{
				"server1": "redis:6379",
			},
		})

		limiter := rate.NewLimiter(ring)
		limiter.Fallback = timerate.NewLimiter(timerate.Every(time.Second), 1000)

		return quotaHandlerWrapper(h, limiter, "service-name")
	}

	handlerFuncSuccess := func(ctx context.Context, req server.Request, rsp interface{}) error {
		return nil
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
			err:         errors.ErrInvalidCtx,
		},
		{
			ctx: micro.NewContext(metadata.NewContext(ctx, map[string]string{
				"Authorization": "",
			}), srv),
			handlerFunc: handlerFuncSuccess,
			req:         request{},
			rsp:         response{},
			err:         errors.ErrNoAuthHeader,
		},
	}

	fn(handlerFuncSuccess)()

}
*/
