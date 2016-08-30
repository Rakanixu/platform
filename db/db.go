package db

import (
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/kazoup/platform/db/srv/engine"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"log"
)

const FileTopic string = "go.micro.topic.files"

func srv(ctx *cli.Context) {

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

	// Init search engine

	if err := engine.Init(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func dbCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run db srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "db",
			Usage:       "DB commands",
			Subcommands: dbCommands(),
		},
	}
}
