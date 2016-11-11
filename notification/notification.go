package notification

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/notification/srv/sockets"
	"github.com/kazoup/platform/notification/srv/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func srv(ctx *cli.Context) {
	service := wrappers.NewKazoupService("notification")

	subscriber.Broker = service.Server().Options().Broker

	// This subscriber receives notification messages and publish same message but over the broker directly
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			subscriber.SubscriberProxy,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Notification handler instantiate with service broker
	// It will allow to subscribe to topics and then stream actions back to clients
	proto.RegisterNotificationHandler(service.Server(), &handler.Notification{
		Server: service.Server(),
	})

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	web := microweb.NewService(microweb.Name("go.micro.web.notification"))

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))

	sockets.NotificationClient = proto.NewNotificationClient(
		"com.kazoup.srv.notification",
		client.DefaultClient,
	)

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
