package main

import (
	"log"

	"github.com/kazoup/platform/bleve/srv/bleve"
	"github.com/kazoup/platform/bleve/srv/handler"
	"github.com/kazoup/platform/bleve/srv/subscriber"
	"github.com/micro/go-micro"
)

const FileTopic string = "go.micro.topic.files"

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.bleve"),
		micro.Version("latest"),
	)
	// Register bleve engine
	engine := bleve.NewBleveEngine()
	//Register SearchBleve Handle
	service.Server().Handle(
		service.Server().NewHandler(
			&handler.SearchBleve{
				Index: enigne.Idx,
			}))

	// Attach indexer subsciber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, &subscriber.FileSubscriber{
			Index: enigine.Idx,
		}),
	); err != nil {
		log.Fatal(err)
	}
	//Init Service
	service.Init()
	//run Servic
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
