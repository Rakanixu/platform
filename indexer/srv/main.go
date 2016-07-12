package main

import (
	"log"
	"time"

	"github.com/kazoup/platform/indexer/srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
)

const topic string = "go.micro.topic.files"

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.indexer"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Init()

	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connert error: %v", err)
	}

	broker.Subscribe(
		topic,
		handler.Subscriber,
		broker.Queue(topic),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
