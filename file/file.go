package file

import (
	"github.com/kazoup/platform/file/srv/handler"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	"log"
)

func srv(ctx *cli.Context) {
	// New service
	service := wrappers.NewKazoupService("file")

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.File{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func fileCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run file srv service",
			Action: srv,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "file",
			Usage:       "File commands",
			Subcommands: fileCommands(),
		},
	}
}
