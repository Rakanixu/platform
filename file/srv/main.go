package main

import (
	"github.com/kazoup/platform/file/srv/handler"
	"github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/file/srv/subscriber"
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
	proto_file.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Subscribers
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-file"),
		),
	)

	// Init DB operations
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	// Init DB bulk indexer
	if err := bulk.Init(service); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
