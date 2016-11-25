package main

import (
	"log"

	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/notification/srv/subscriber"
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
	if err := service.Server().Handle(
		service.Server().NewHandler(
			&handler.Notification{
				Server: service.Server(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	service.Init()
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
