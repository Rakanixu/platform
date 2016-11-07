package main

import (
	"github.com/kazoup/platform/file/srv/handler"
	"github.com/kazoup/platform/lib/wrappers"
	"log"
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
