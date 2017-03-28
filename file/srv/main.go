package main

import (
	"github.com/kazoup/platform/file/srv/handler"
	"github.com/kazoup/platform/file/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
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

	// Healthchecks for file-srv
	healthchecks.RegisterFileHealthChecks(service, m)
	healthchecks.RegisterBrokerHealthChecks(service, m)

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.File{
			Client: service.Client(),
		}),
	)

	// Subscribers
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			&subscriber.AnnounceFile{
				Client: service.Client(),
				Broker: service.Server().Options().Broker,
			},
			server.SubscriberQueue("announce-file"),
		),
	)

	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
