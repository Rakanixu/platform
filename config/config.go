package config

import (
	"github.com/kazoup/platform/config/web/handler"
	"github.com/kazoup/platform/config/web/sockets"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"time"
)

func web(ctx *cli.Context) {
	var m monitor.Monitor

	web := microweb.NewService(microweb.Name("go.micro.web.config"))

	// config-web monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterConfigWebHealthChecks(m)

	// Attach web handler (socket)
	web.Handle("/platform/ping", websocket.Handler(sockets.PingPlatform))
	web.HandleFunc("/health", handler.HandleHealthCheck)
	web.Run()
}

func configCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "web",
			Usage:  "Run config web service",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "config",
			Usage:       "Auth commands",
			Subcommands: configCommands(),
		},
	}
}
