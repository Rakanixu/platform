package main

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"log"
)

func main() {
	// New service

	service := wrappers.NewKazoupService("datasource")

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, handler.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.DataSource{
			Client: service.Client(),
		}),
	)

	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
