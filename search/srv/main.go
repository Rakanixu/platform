package main

import (
	"github.com/kazoup/platform/search/srv/engine"
	"log"
	//_ "github.com/kazoup/platform/search/srv/engine/bleve"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/go-micro"
)

const (
	elasticServiceName = "go.micro.srv.db"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.search"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			ElasticServiceName: elasticServiceName,
			Client:             service.Client(),
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
