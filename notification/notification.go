package notification

import (
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	"github.com/kazoup/platform/notification/srv/subscriber"
	"github.com/kazoup/platform/notification/web/sockets"
	"github.com/micro/cli"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("notification")

	// This subscriber receives notification messages and publish same message but over the broker directly
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			&subscriber.Proxy{
				Broker: service.Server().Options().Broker,
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	// Notification handler instantiate with service server
	// It will allow to subscribe to topics and then stream actions back to clients
	if err := service.Server().Handle(
		service.Server().NewHandler(
			&handler.Notification{
				Server: service.Server(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	web := microweb.NewService(microweb.Name("go.micro.web.notification"))

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))

	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}

func notificationCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run notification srv service",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run notification web service",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "notification",
			Usage:       "Notification commands",
			Subcommands: notificationCommands(),
		},
	}
}
