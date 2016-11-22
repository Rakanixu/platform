package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"golang.org/x/net/context"
)

var Broker broker.Broker

// SubscriberProxy listens for messages and proxys to service Broker to be streamed to clients afterwards
func SubscriberProxy(ctx context.Context, notificationMsg *proto.NotificationMessage) error {
	b, err := json.Marshal(notificationMsg)
	if err != nil {
		return err
	}
	// Publish on the broker, it allows to handle data properly in broker Handler
	if err := Broker.Publish(globals.NotificationTopic, &broker.Message{
		Body: b,
	}); err != nil {
		return err
	}

	return nil
}