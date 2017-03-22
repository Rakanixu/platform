package main

import (
	"github.com/kazoup/platform/audioenrich/srv/handler"
	"github.com/kazoup/platform/audioenrich/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
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

	service := wrappers.NewKazoupService("audioenrich", m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach handler
	if err := service.Server().Handle(
		service.Server().NewHandler(
			&handler.AudioEnrich{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	gcslib.Register()

	s := &subscriber.Enrich{
		Client:             service.Client(),
		GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		EnrichMsgChan:      make(chan subscriber.EnrichMsgChan, 1000000),
		Workers:            20,
	}
	subscriber.StartWorkers(s)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AudioEnrichTopic,
			s,
			server.SubscriberQueue("audioenrich"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			&subscriber.AnnounceAudioEnrich{
				Client: service.Client(),
				Broker: service.Server().Options().Broker,
			},
			server.SubscriberQueue("announce-audioenrich"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Init service
	service.Init()

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
