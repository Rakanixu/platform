package auth

import (
	"log"
	"time"

	api_handler "github.com/kazoup/platform/config/api/handler"

	srv_handler "github.com/kazoup/platform/config/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func api(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.api.config"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(api_handler.Config)),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(srv_handler.Config)),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
func configCommands() []cli.Command {
	return []cli.Command{{
		Name:   "api",
		Usage:  "Run config api service",
		Action: api,
	}, {
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
