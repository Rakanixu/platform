package main

import (
	"github.com/kazoup/platform/flag/srv/handler"
	"github.com/micro/go-micro"
	"log"
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
		service.Server().NewHandler(&handler.Flag{
			DbServiceName: "go.micro.srv.db",
			Client:        service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
