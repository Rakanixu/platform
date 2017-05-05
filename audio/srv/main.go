package main

import (
	"github.com/kazoup/platform/audio/srv/handler"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	"github.com/kazoup/platform/audio/srv/subscriber"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	"github.com/kazoup/platform/lib/objectstorage"
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

	if err := custom.Init(); err != nil {
		log.Fatal(err)
	}

	var m monitor.Monitor

	service := wrappers.NewKazoupService("audio", m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach handler
	proto_audio.RegisterServiceHandler(service.Server(), new(handler.Service))

	objectstorage.Init()

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AudioEnrichTopic,
			subscriber.NewTaskHandler(20),
			server.SubscriberQueue("audio"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-audio"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
