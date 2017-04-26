package main

import (
	"github.com/kazoup/platform/document/srv/handler"
	"github.com/kazoup/platform/document/srv/proto/document"
	"github.com/kazoup/platform/document/srv/subscriber"
	"github.com/kazoup/platform/lib/db/custom"
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
	// Init DB Operations
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	// Init DB Custom Operations
	if err := custom.Init(); err != nil {
		log.Fatal(err)
	}

	var m monitor.Monitor

	service := wrappers.NewKazoupService("document", m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach handler
	proto_document.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DocEnrichTopic,
			subscriber.NewTaskHandler(25),
			server.SubscriberQueue("document"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-document"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
