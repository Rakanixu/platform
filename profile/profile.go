package quota

import (
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/profile/srv/handler"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	// New service
	service := wrappers.NewKazoupService("profile", m)

	// Monitor for file-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Healthchecks for quota-srv
	healthchecks.RegisterBrokerHealthChecks(service, m)

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Profile{
			Client: service.Client(),
		}),
	)

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Subscription{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func profileCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run profile srv service",
			Action: srv,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "profile",
			Usage:       "Profile commands",
			Subcommands: profileCommands(),
		},
	}
}
