package main

import (
	"log"

	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/notification/srv/subscriber"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	service := wrappers.NewKazoupService("notification")

	// This subscriber receives notification messages and publish same message but over the broker directly
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			&subscriber.Proxy{
				Broker: service.Server().Options().Broker,
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	// Notification handler instantiate with service broker
	// It will allow to subscribe to topics and then stream actions back to clients
	proto.RegisterNotificationHandler(service.Server(), &handler.Notification{
		Server: service.Server(),
	})

	service.Init()
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
