package main

import (
	"log"

	"github.com/kazoup/platform/file/srv/handler"
	"github.com/kazoup/platform/lib/wrappers"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	// New service
	service := wrappers.NewKazoupService("file")

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.File{
			Client: service.Client(),
		}),
	)

	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
