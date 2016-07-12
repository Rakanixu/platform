package main

import (
	"log"

	"github.com/kazoup/platform/flag/srv/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.flag"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Flag)),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
