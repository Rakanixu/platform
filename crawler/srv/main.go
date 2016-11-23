package main

import (
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	"log"
)

func main() {
	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	service := wrappers.NewKazoupService("crawler")

	// Init srv
	service.Init()

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.ScanTopic,
			&subscriber.Crawler{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
