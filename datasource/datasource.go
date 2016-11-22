package datasource

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("datasource")

	// Init broker on subscriber
	// This is required to be able to handle the data properly when
	// we want to stream the messages over the notification socket
	subscriber.Broker = service.Server().Options().Broker

	// Attach crawler started subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerStartedTopic, subscriber.SubscribeCrawlerStarted)); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, subscriber.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// Attach delete bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteBucketTopic, subscriber.SubscribeDeleteBucket)); err != nil {
		log.Fatal(err)
	}

	// Attach clean bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteFileInBucketTopic, subscriber.SubscribeDeleteFileInBucket)); err != nil {
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
