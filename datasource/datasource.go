package datasource

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/cli"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("datasource")

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, handler.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.DataSource{
			Client: service.Client(),
		}),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func datasourceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run datasource srv service",
			Action: srv,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "datasource",
			Usage:       "Datasource commands",
			Subcommands: datasourceCommands(),
		},
	}
}
