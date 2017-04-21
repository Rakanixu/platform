package main

import (
	"github.com/kazoup/platform/entities/srv/handler"
	"github.com/kazoup/platform/entities/srv/proto/entities"
	"github.com/kazoup/platform/entities/srv/subscriber"
	"github.com/kazoup/platform/lib/db/operations"
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
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	var m monitor.Monitor

	service := wrappers.NewKazoupService("entities", m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach handler
	proto_entities.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.ExtractEntitiesTopic,
			subscriber.NewTaskHandler(40),
			server.SubscriberQueue("entities"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-entities"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
