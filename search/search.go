package elastic

import (
	"github.com/kazoup/platform/search/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
)

func srv(ctx *cli.Context) {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.search"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Search{
			Client:             service.Client(),
			ElasticServiceName: "go.micro.srv.elastic",
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func elasticCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run search srv",
			Action: srv,
		},
	}
}
