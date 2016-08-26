package main

import (
	"github.com/kazoup/platform/config/srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"log"
)

//go-bindata -o data/bindata.go -pkg data data
func main() {
	cmd.Init()

	// New service
	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.Version("latest"),
	)

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Config{
			Client:        service.Client(),
			DbServiceName: "go.micro.srv.flag",
		}),
	)

	// Initialize service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
