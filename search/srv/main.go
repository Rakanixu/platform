package main

import (
	"github.com/kazoup/platform/search/srv/engine"
	"log"
	//_ "github.com/kazoup/platform/search/srv/engine/bleve"
	"github.com/kazoup/platform/lib/wrappers"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	"github.com/kazoup/platform/search/srv/handler"
	_ "github.com/micro/go-plugins/broker/nats"
)

func main() {
	// New Service
	service := wrappers.NewKazoupService("search")
	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			Client: service.Client(),
		}),
	)

	// Initialise service
	service.Init()
	// Init search engine

	if err := engine.Init(); err != nil {
		log.Fatal(err)
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
