package notification

import (
	"github.com/kazoup/platform/notification/srv/sockets"
	"github.com/kazoup/platform/notification/srv/subscriber"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/wrappers"
	"github.com/micro/cli"
	_ "github.com/micro/go-plugins/broker/nats"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("notification")

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			subscriber.Notify,
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

	// Attach web handler (socket)
	web.Handle("/platform/notify", websocket.Handler(sockets.Notify))
	web.Run()
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
