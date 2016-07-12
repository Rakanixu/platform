package main

import (
	"log"

	"github.com/kazoup/platform/crawler/srv/handler"
	proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/broker/nats"
)

func main() {
	cmd.Init()

	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connert error: %v", err)
	}

	service := micro.NewService(
		// TODO: com.kazoup.srv.crawler
		micro.Name("go.micro.srv.crawler"),
		micro.Version("latest"),
	)

	service.Init()

	proto.RegisterCrawlHandler(service.Server(), new(handler.Crawl))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
