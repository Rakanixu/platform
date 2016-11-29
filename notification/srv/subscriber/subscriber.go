package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
	"golang.org/x/net/context"
)

type Proxy struct {
	Client client.Client
	Server server.Server
}

// SubscriberProxy listens for messages and proxys to service Broker to be streamed to clients afterwards
func (p *Proxy) SubscriberProxy(ctx context.Context, notificationMsg *proto.NotificationMessage) error {
	b, err := json.Marshal(notificationMsg)
	if err != nil {
		return err
	}

	fmt.Println("ADDRESS SUBSPROXY 1", p.Server.Options().Broker.Address())
	fmt.Println("ADDRESS SUBSPROXY 2", p.Client.Options().Broker.Address())

	// Publish on the broker, it allows to handle data properly in broker Handler
	if err := p.Server.Options().Broker.Publish(globals.NotificationProxyTopic, &broker.Message{
		Body: b,
	}); err != nil {
		return err
	}

	return nil
}
