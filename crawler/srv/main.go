package main

import (
	"log"

	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
)

func main() {
	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	service := wrappers.NewKazoupService("crawler")

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
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
