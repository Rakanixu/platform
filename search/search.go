package search

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/search/srv/engine"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	// New Service
	service := wrappers.NewKazoupService("search", globals.QUOTA_HANDLER_SEARCH, globals.QUOTA_SUBS_SEARCH, m)

	// Instantiate monitor for search-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterSearchHealthChecks(service, m)
	healthchecks.RegisterBrokerHealthChecks(service, m)

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
