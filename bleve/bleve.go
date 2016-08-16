package config

import (
	"log"

	"github.com/kazoup/platform/bleve/srv/bleve"
	"github.com/kazoup/platform/bleve/srv/handler"
	"github.com/kazoup/platform/bleve/srv/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func srv(ctx *cli.Context) {

	service := micro.NewService(
		micro.Name("go.micro.srv.bleve"),
		micro.Version("latest"),
	)
	// Register bleve engine
	engine := bleve.NewBleveEngine()
	//Register SearchBleve Handle
	service.Server().Handle(
		service.Server().NewHandler(
			&handler.SearchBleve{
				Index: engine.Idx,
			}))
	test
	// Attach indexer subsciber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, &subscriber.FileSubscriber{
			MsgCh: engine.FileChannel,
		},
		)); err != nil {
		log.Fatal(err)
	}
	// Init Bacth indexer
	engine.Indexer()

	//run Servic
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func bleveCommands() []cli.Command {
	return []cli.Command{
		Name:   "srv",
		Usage:  "Run bleve  service",
		Action: srv,
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "bleve",
			Usage:       "Bleve commands",
			Subcommands: bleveCommands(),
		},
	}
}
