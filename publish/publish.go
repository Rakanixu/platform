package publish

import (
	"log"

	"github.com/kazoup/platform/publish/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
)

func srv(ctx *cli.Context) {

	// Init broker
	if err := broker.Init(); err != nil {
		log.Fatalf("Broker Init error: %v", err)
	}

	// Connect broker
	if err := broker.Connect(); err != nil {
		log.Fatalf("Broker Connert error: %v", err)
	}

	// New Service
	service := micro.NewService(
		//TODO: com.kazoup.srv.publisher
		micro.Name("go.micro.srv.publish"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.Publish)),
	)
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func publishCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run publish service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "publish",
			Usage:       "Publish commands",
			Subcommands: publishCommands(),
		},
	}
}
