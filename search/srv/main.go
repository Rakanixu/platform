package main

import (
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/go-micro"
	"log"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.search"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			Client:             service.Client(),
			ElasticServiceName: "go.micro.srv.elastic",
		}),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
