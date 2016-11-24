package main

import (
	"log"

	"github.com/kazoup/platform/config/srv/handler"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro/cmd"
)

//go-bindata -o data/bindata.go -pkg data data
func main() {
	cmd.Init()

	// New service

	service := wrappers.NewKazoupService("config")

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Config{
			Client: service.Client(),
		}),
	)

	// Initialize service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
