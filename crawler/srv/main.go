package main

import (
	"github.com/kazoup/platform/crawler/srv/subscriber"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/db/bulk"
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
	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	var m monitor.Monitor

	service := wrappers.NewKazoupService("crawler", m)

	// crawler-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-crawler"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DiscoverTopic,
			new(subscriber.TaskHandler),
			server.SubscriberQueue("crawler"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Initialize DB operations
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	// Initialize DB bulk indexer
	if err := bulk.Init(service); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
