package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
)

func StreamNotifications(n *Notification, req *proto.StreamRequest) (chan *proto.NotificationMessage, chan bool, error) {
	che := make(chan *proto.NotificationMessage, 10000) // To be sure channel is not blocked
	exit := make(chan bool)

	// We subscribe directly to the broker to be able to handle the data internally
	sub, err := n.Server.Options().Broker.Subscribe(globals.NotificationProxyTopic, func(p broker.Publication) error {
		var e *proto.NotificationMessage

		if err := json.Unmarshal(p.Message().Body, &e); err != nil {
			return err
		}

		if req.UserId == e.UserId {
			che <- e
		}

		return nil
	})

	if err != nil {
		fmt.Println("ERROR attaching subscriber NotificationProxyTopic", err)
		return nil, nil, err
	}

	go func() {
		<-exit
		sub.Unsubscribe()
	}()

	return che, exit, nil
}
