package auth

import (
	"log"
	"time"

	"github.com/kazoup/platform/crawler/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/kazoup/platform/crawler/srv/subscriber"
)

const topic string = "go.micro.topic.scan"
func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.crawler"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(handler.Crawl)),
	)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			topic,
			subscriber.Scans,
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func crawlerCommands() []cli.Command {
	return []cli.Command{{
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
