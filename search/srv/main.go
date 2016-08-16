package main

import (
	"log"

	"github.com/kazoup/platform/search/srv/engine"

	_ "github.com/kazoup/platform/search/srv/engine/bleve"
	//_ "github.com/kazoup/platform/search/srv/engine/elastic"
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.search"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{}),
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
