package auth

import (
	"log"
	"time"

	"github.com/kazoup/platform/auth/api/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func api(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.api.auth"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(handler.Auth)),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func authCommands() []cli.Command {
	return []cli.Command{{
		Name:   "api",
		Usage:  "Run auth api service",
		Action: api,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "auth",
			Usage:       "Auth commands",
			Subcommands: authCommands(),
		},
	}
}
