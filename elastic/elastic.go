package elastic

import (
	"log"
	"strings"

	elastic_api "github.com/kazoup/platform/elastic/api/handler"
	elastic "github.com/kazoup/platform/elastic/srv/elastic"
	elastic_srv "github.com/kazoup/platform/elastic/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func api(ctx *cli.Context) {

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.elastic"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(elastic_api.Elastic)),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func srv(ctx *cli.Context) {

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.elastic"),
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
			parts := strings.Split(c.String("elasticsearch_hosts"), ",")
			elastic.Hosts = parts
		}),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(elastic_srv.Elastic)),
	)

	// Initialise elasticsearch
	elastic.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func elasticCommands() []cli.Command {
	return []cli.Command{{
		Name:   "api",
		Usage:  "Run elastic api service",
		Action: api,
	}, {

		Name:   "srv",
		Usage:  "Run elastic srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "elastic",
			Usage:       "Elastic commands",
			Subcommands: elasticCommands(),
		},
	}
}
