package main

import (
	"github.com/kazoup/platform/db/srv/engine"
	//_ "github.com/kazoup/platform/db/srv/engine/bleve"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
)

const FileTopic string = "go.micro.topic.files"

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.db"),
		micro.Version("latest"),
		micro.Flags(
			cli.StringFlag{
				Name:   "elasticsearch_hosts",
				EnvVar: "ELASTICSEARCH_HOSTS",
				Usage:  "Comma separated list of elasticsearch hosts",
				Value:  "localhost:9200",
			},
		),
		micro.Action(func(c *cli.Context) {
			//parts := strings.Split(c.String("elasticsearch_hosts"), ",")
			//elastic.Hosts = parts
		}),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DB)),
	)

	// Attach indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, engine.Subscribe)); err != nil {
		log.Fatal(err)
	}

	// Initialise service
	service.Init()
	// Init search engine

	if err := engine.Init(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
