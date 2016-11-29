package handler

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/nats"
)

func StreamNotifications(n *Notification, req *proto.StreamRequest) (chan *proto.NotificationMessage, chan bool, error) {
	che := make(chan *proto.NotificationMessage, 10000) // To be sure channel is not blocked
	exit := make(chan bool)

	fmt.Println("ADDRESS 1", broker.DefaultBroker.Address())
	fmt.Println("ADDRESS 2", client.DefaultClient.Options().Broker.Address())
	fmt.Println("ADDRESS 3", n.Client.Options().Broker.Address())
	fmt.Println("ADDRESS 4", n.Server.Options().Broker.Address())
	fmt.Println("ADDRESS 5", cmd.DefaultBrokers["nats"])

	// We subscribe directly to the broker to be able to handle the data internally
	sub, err := broker.DefaultBroker.Subscribe(globals.NotificationProxyTopic, func(p broker.Publication) error {
		var e *proto.NotificationMessage

		if err := json.Unmarshal(p.Message().Body, &e); err != nil {
			return err
		}

		fmt.Println("s.Options().Broker.Subscribe", req, e)

		if req.UserId == e.UserId {
			che <- e
		}

		return nil
	})

	if err != nil {
		fmt.Println("ERROR s.Options().Broker.Subscribe", err)
		return nil, nil, err
	}

	go func() {
		<-exit
		sub.Unsubscribe()
	}()

	fmt.Println("Channels returned")

	return che, exit, nil
}
