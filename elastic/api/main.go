package main

import (
	"log"

	"github.com/kazoup/platform/elastic/api/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.elastic"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Elastic)),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
