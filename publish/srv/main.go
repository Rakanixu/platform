package main

import (
	"log"

	"github.com/kazoup/platform/publish/srv/handler"
	publish "github.com/kazoup/platform/publish/srv/proto/publish"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
)

func main() {
	// Init broker
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	// Connect broker
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connert error: %v", err)
	}

	// New Service
	service := micro.NewService(
		//TODO: com.kazoup.srv.publisher
		micro.Name("go.micro.srv.publish"),
		micro.Version("latest"),
	)

	// Register Handler
	publish.RegisterPublishHandler(service.Server(), new(handler.Publish))

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
