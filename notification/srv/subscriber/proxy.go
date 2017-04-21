package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"golang.org/x/net/context"
)

type ProxyHandler struct{}

// SubscriberProxy listens for messages and proxys to service Broker to be streamed to clients afterwards
func (p *ProxyHandler) SubscriberProxy(ctx context.Context, notificationMsg *proto_notification.NotificationMessage) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

	b, err := json.Marshal(notificationMsg)
	if err != nil {
		return err
	}

	// Publish on the broker, it allows to handle data properly in broker Handler
	if err := srv.Server().Options().Broker.Publish(globals.NotificationProxyTopic, &broker.Message{
		Body: b,
	}); err != nil {
		return err
	}

	return nil
}
