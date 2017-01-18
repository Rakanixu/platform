package db

import (
	"github.com/kazoup/platform/db/srv/engine"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	// New Service
	service := wrappers.NewKazoupService("db", m)

	m = monitor.NewMonitor(
		monitor.Interval(time.Second),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DB)),
	)

	// Register Config Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Config)),
	)

	// Attach file indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.FilesTopic, &engine.Files{
			Client: service.Client(),
		})); err != nil {
		log.Fatal(err)
	}

	// Attach slack user indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.SlackUsersTopic, engine.SubscribeSlackUsers)); err != nil {
		log.Fatal(err)
	}

	// Attach slack channel indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.SlackChannelsTopic, engine.SubscribeSlackChannels)); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, engine.SubscribeCrawlerFinished)); err != nil {
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
