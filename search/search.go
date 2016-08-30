package search

import (
	"log"

	"github.com/kazoup/platform/search/srv/engine"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

const (
	elasticServiceName = "go.micro.srv.db"
)

func srv(ctx *cli.Context) {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.search"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			ElasticServiceName: elasticServiceName,
			Client:             service.Client(),
		}),
	)

	if err := engine.Init(); err != nil {
		log.Fatal(err)
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func searchCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run search srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "config",
			Usage:       "Search commands",
			Subcommands: searchCommands(),
		},
	}
}
