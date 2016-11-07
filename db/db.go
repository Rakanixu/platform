package db

import (
	"github.com/kazoup/platform/db/srv/engine"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

func srv(ctx *cli.Context) {
	// New Service
	service := wrappers.NewKazoupService("db")

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DB)),
	)

	// Attach file indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.FilesTopic, engine.SubscribeFiles)); err != nil {
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
