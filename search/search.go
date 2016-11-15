package search

import (
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/search/srv/engine"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	"log"
)

func srv(ctx *cli.Context) {
	// New Service
	service := wrappers.NewKazoupService("search")
	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			Client: service.Client(),
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
			Name:        "search",
			Usage:       "Search commands",
			Subcommands: searchCommands(),
		},
	}
}
