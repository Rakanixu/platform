package config

import (
	srv_handler "github.com/kazoup/platform/config/srv/handler"
	"github.com/kazoup/platform/config/srv/sockets"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"

	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("config")

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&srv_handler.Config{
			Client: service.Client(),
		}),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func web(ctx *cli.Context) {
	web := microweb.NewService(microweb.Name("go.micro.web.config"))

	// Attach web handler (socket)
	web.Handle("/platform/ping", websocket.Handler(sockets.PingPlatform))
	web.Run()
}

func configCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run config srv service",
			Action: srv,
		},
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
