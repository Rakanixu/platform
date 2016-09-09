package main

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro"
	"log"
)

func main() {
	// New service
	service := micro.NewService(
		micro.Name("go.micro.srv.datasource"),
		micro.Version("latest"),
	)

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, handler.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.DataSource{
			Client:             service.Client(),
			ElasticServiceName: "go.micro.srv.db",
		}),
	)

	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
