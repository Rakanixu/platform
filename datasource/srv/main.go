package main

import (
	"log"

	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/micro/go-micro"
)

func main() {
	// New service
	service := micro.NewService(
		micro.Name("go.micro.srv.datasource"),
		micro.Version("latest"),
	)

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DataSource)),
	)

	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
