package datasource

import (
	"log"
	"time"

	datasource_api "github.com/kazoup/platform/datasource/api/handler"
	datasource_srv "github.com/kazoup/platform/datasource/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func api(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.api.datasource"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(datasource_api.Datasource)),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.datasource"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),
	)

	service.Server().Handle(
		service.Server().NewHandler(new(datasource_srv.DataSource)),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
func datasourceCommands() []cli.Command {
	return []cli.Command{{
		Name:   "api",
		Usage:  "Run datasource api service",
		Action: api,
	}, {

		Name:   "srv",
		Usage:  "Run datasource srv service",
		Action: srv,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "datasource",
			Usage:       "Datasource commands",
			Subcommands: datasourceCommands(),
		},
	}
}
