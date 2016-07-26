package indexer

import (
	"log"
	"time"

	"github.com/kazoup/platform/indexer/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/nats"
)

const topic string = "go.micro.topic.files"

func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.indexer"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(topic, handler.Subscriber),
	); err != nil {
		log.Printf("Got error subscibing .. %s", err.Error())
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func indexerCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run indexer service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "indexer",
			Usage:       "Indexer commands",
			Subcommands: indexerCommands(),
		},
	}
}
