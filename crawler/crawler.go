package crawler

import (
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/lib/categories"
	//"github.com/kazoup/platform/lib/globals"
	"log"

	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("crawler")

	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.ScanTopic,
			&subscriber.Crawler{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func crawlerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run crawler service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "crawler",
			Usage:       "Crawler commands",
			Subcommands: crawlerCommands(),
		},
	}
}
