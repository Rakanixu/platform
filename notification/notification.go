package notification

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	"github.com/kazoup/platform/notification/srv/subscriber"
	web_handler "github.com/kazoup/platform/notification/web/handler"
	"github.com/kazoup/platform/notification/web/sockets"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("notification", m)

	// Monitor for notification-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterNotificationSrvHealthChecks(service, m)
	healthchecks.RegisterBrokerHealthChecks(service, m)

	// This subscriber receives notification messages and publish same message but over the broker directly
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.NotificationTopic,
			&subscriber.Proxy{
				Client: service.Client(),
				Server: service.Server(),
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
				Client: service.Client(),
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
	var m monitor.Monitor

	web := microweb.NewService(
		microweb.Name("go.micro.web.notification"),
	)

	// Monitor for notification-web
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterNotificationWebHealthChecks(m)

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))
	web.HandleFunc("/health", web_handler.HandleHealthCheck)

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
