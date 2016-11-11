package main

import (
	"github.com/kazoup/platform/flag/srv/handler"
	"github.com/kazoup/platform/lib/wrappers"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

func main() {
	// New Service
	service := wrappers.NewKazoupService("flag")

	// Initialise service
	service.Init()

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Flag{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
