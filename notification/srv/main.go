package main

import (
	"github.com/kazoup/platform/notification/srv/subscriber"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

func main() {
	service := wrappers.NewKazoupService("notification")

	// Init srv
	service.Init()

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			subscriber.Notify,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
