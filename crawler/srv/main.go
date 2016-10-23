package main

import (
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/structs/categories"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
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
			subscriber.Scans,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
