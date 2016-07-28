package main

import (
	"log"
	"time"

	"github.com/kazoup/platform/indexer/srv/handler"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
)

const topic string = "go.micro.topic.files"

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.indexer"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(topic, handler.FileSubscriber),
	); err != nil {
		log.Printf("Got error subscibing .. %s", err.Error())
	}

	service.Init()
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
