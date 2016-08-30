package flag

import (
	"github.com/kazoup/platform/flag/srv/handler"
	"github.com/micro/go-micro"

	"github.com/micro/cli"
	"log"
)

func srv(ctx *cli.Context) {

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.flag"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Flag{
			DbServiceName: "go.micro.srv.db",
			Client:        service.Client(),
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
