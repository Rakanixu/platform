package main

import (
	"github.com/kazoup/platform/crawler/srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

const topic string = "go.micro.topic.scan"

func main() {
	cmd.Init()

	service := micro.NewService(
		// TODO: com.kazoup.srv.crawler
		micro.Name("go.micro.srv.crawler"),
		micro.Version("latest"),
	)

	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(
			new(handler.Crawl),
		),
	)

	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			topic,
			handler.Subscriber,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
