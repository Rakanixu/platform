package main

import (
	"log"

	"github.com/kazoup/platform/bleve/srv/bleve"
	"github.com/kazoup/platform/bleve/srv/handler"
	"github.com/kazoup/platform/bleve/srv/subscriber"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
)

//FileTopic topic to scan files
const FileTopic string = "go.micro.topic.files"

func main() {

	service := micro.NewService(
		micro.Name("go.micro.srv.bleve"),
		micro.Version("latest"),
	)
	// Register bleve engine
	files := bleve.NewBleveEngine("/tmp/files.idx")
	//Register SearchBleve Handle
	service.Server().Handle(
		service.Server().NewHandler(
			&handler.SearchBleve{
				Index: files.Idx,
			}))

	// Attach indexer subsciber

	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, &subscriber.FileSubscriber{
			MsgCh: files.FileChannel,
		},
		)); err != nil {
		log.Fatal(err)
	}
	// Init Bacth indexer
	files.Indexer()

	//Init Service
	service.Init()
	//run Servic
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
