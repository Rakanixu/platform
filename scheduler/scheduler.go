package scheduler

import (
	"log"

	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/scheduler/srv/handler"
	"github.com/micro/cli"
)

func srv(ctx *cli.Context) {
	// New Service
	service := wrappers.NewKazoupService("scheduler")

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Scheduler{
			Client: service.Client(),
			Crons:  make([]*handler.CronWrapper, 0),
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
