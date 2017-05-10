package subscriber

import (
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"testing"
)

var (
	proxyHandler = new(ProxyHandler)
)

func TestSubscriberProxy(t *testing.T) {
	var subscriberProxyTestData = []struct {
		ctx context.Context
		msg *proto_notification.NotificationMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&proto_notification.NotificationMessage{
				Info: "test",
			},
			nil,
		},
		//Invalid context
		{
			ctx,
			&proto_notification.NotificationMessage{
				Info: "test",
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range subscriberProxyTestData {
		err := proxyHandler.SubscriberProxy(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}
