package config

import (
	"log"
	"time"

	srv_handler "github.com/kazoup/platform/config/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)


func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&srv_handler.Config{
			Client:        service.Client(),
			DbServiceName: "go.micro.srv.flag",
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
