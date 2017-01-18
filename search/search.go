package search

import (
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

type request struct {
	service string
	method  string
}

func (r *request) Service() string {
	return r.service
}

func (r *request) Method() string {
	return r.method
}

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	// New Service
	service := wrappers.NewKazoupService("search", m)

	// Instantiate monitor for search-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Second),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

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
