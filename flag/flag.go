package flag

import (
	"github.com/kazoup/platform/flag/srv/handler"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

func srv(ctx *cli.Context) {

	// New Service
	service := wrappers.NewKazoupService("flag")

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Flag{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func flagCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run flag srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "flag",
			Usage:       "Flag commands",
			Subcommands: flagCommands(),
		},
	}
}
