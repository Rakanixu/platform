package main

import (
	"io/ioutil"
	"log"

	"github.com/kazoup/platform/crawler/srv/handler"
	proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/structs/categories"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
)

const topic string = "go.micro.topic.scan"

func main() {
	// Load categories JSON map. categories_map.json
	mapping, err := ioutil.ReadFile("categories_map.json")
	if err != nil {
		log.Fatal(err)
	}

	if err := categories.SetMap(mapping); err != nil {
		log.Fatal(err)
	}

	service := micro.NewService(
		micro.Name("go.micro.srv.crawler"),
		micro.Version("latest"),
	)

	// Init srv
	service.Init()
	proto.RegisterCrawlHandler(service.Server(), new(handler.Crawl))
	// Attach handler
	//service.Server().Handle(
	//	service.Server().NewHandler(
	//		new(handler.Crawl),
	//	),
	//)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			topic,
			subscriber.Scans,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
