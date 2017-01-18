package datasource

import (
	"log"

	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("datasource", m)

	// datasource-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Attach crawler started subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerStartedTopic, &subscriber.CrawlerStarted{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, &subscriber.CrawlerFinished{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach delete bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteBucketTopic, &subscriber.DeleteBucket{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach clean bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteFileInBucketTopic, &subscriber.DeleteFileInBucket{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
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
