package main

import (
	"log"

	data "github.com/kazoup/platform/config/srv/data"
	"github.com/kazoup/platform/config/srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
)

//go-bindata -o data/bindata.go -pkg data data
func main() {
	cmd.Init()
	es_flags, err := data.Asset("data/es_flags.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}
	es_mapping, err := data.Asset("data/es_mapping_files.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}
	es_settings, err := data.Asset("data/es_settings.json")
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}
	// New service
	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.Version("latest"),
	)

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Config{
			Client:             service.Client(),
			ESSettings:         &es_settings,
			ESFlags:            &es_flags,
			ESMapping:          &es_mapping,
			ElasticServiceName: "go.micro.srv.elastic",
		}),
	)

	// Initialize service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
