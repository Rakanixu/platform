package main

import (
	"log"

	"github.com/kazoup/platform/policy/srv/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.policy"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Policy)),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
