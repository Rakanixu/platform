package main

import (
	"log"
	"time"

	"github.com/kazoup/platform/auth/api/handler"
	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		//TODO change namespace to com.kazoup.api
		micro.Name("go.micro.api.auth"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Init()

	service.Server().Handle(
		service.Server().NewHandler(new(handler.Auth)),
	)

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
