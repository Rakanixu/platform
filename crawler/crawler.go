package crawler

import (
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	var m monitor.Monitor

	service := wrappers.NewKazoupService("crawler", globals.QUOTA_HANDLER_CRAWLER, globals.QUOTA_SUBS_CRAWLER, m)

	// crawler-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

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
