package main

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/notification/srv/handler"
	"github.com/kazoup/platform/notification/srv/subscriber"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
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
				Server: service.Server(),
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	// React to tasks done, notify user about them
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			&subscriber.AnnounceNotification{
				Client: service.Client(),
				Broker: service.Server().Options().Broker,
			},
			server.SubscriberQueue("announce-notification"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Notification handler instantiate with service broker
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

	service.Init()
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
