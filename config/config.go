package config

import (
	srv_handler "github.com/kazoup/platform/config/srv/handler"
	"github.com/kazoup/platform/config/srv/sockets"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/cli"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("config")
	web := microweb.NewService(microweb.Name("go.micro.web.config"))

	// Attach handler
	service.Server().Handle(
		service.Server().NewHandler(&srv_handler.Config{
			Client: service.Client(),
		}),
	)
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}

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
