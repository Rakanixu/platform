package file

import (
	"log"

	"github.com/kazoup/platform/file/srv/handler"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	// New service
	service := wrappers.NewKazoupService("file", m)

	// Monitor for file-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

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
