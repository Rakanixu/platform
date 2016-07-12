package main

import (
	"log"

	"github.com/kazoup/platform/config/srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
)

func main() {
	cmd.Init()

	// New service
	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.Version("latest"),
	)

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Config)),
	)

	// Initialize service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
