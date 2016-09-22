package config

import (
	"log"

	srv_handler "github.com/kazoup/platform/config/srv/handler"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/cli"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("config")

	// Attach handler

	service.Server().Handle(
		service.Server().NewHandler(&srv_handler.Config{
			Client: service.Client(),
		}),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
func configCommands() []cli.Command {
	return []cli.Command{{
		Name:   "srv",
		Usage:  "Run config srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "config",
			Usage:       "Auth commands",
			Subcommands: configCommands(),
		},
	}
}
