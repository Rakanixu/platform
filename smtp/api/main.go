package main

import (
	"log"

	"github.com/kazoup/smtp/api/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.smtp"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Smtp)),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
