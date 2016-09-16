package scheduler

import (
	"github.com/kazoup/platform/scheduler/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"log"
)

func srv(ctx *cli.Context) {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.scheduler"),
		micro.Version("latest"),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Scheduler{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
func schedulerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run sheduler srv service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "scheduler",
			Usage:       "Scheduler commands",
			Subcommands: schedulerCommands(),
		},
	}
}
