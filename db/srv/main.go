package main

import (
	"github.com/kazoup/platform/db/srv/engine"
	_ "github.com/kazoup/platform/db/srv/engine/bleve"
	"log"
	//_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"github.com/micro/go-micro"
)

const FileTopic string = "go.micro.topic.files"

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.db"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DB)),
	)

	// Attach indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, engine.Subscribe)); err != nil {
		log.Fatal(err)
	}

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
