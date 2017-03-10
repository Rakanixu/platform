package main

import (
	"github.com/kazoup/platform/docenrich/srv/handler"
	"github.com/kazoup/platform/docenrich/srv/subscriber"
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

	service := wrappers.NewKazoupService("docenrich", m)

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
			&handler.DocEnrich{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Fatal(err)
	}

	s := &subscriber.Enrich{
		Client:        service.Client(),
		EnrichMsgChan: make(chan subscriber.EnrichMsgChan, 1000000),
		Workers:       25,
	}
	subscriber.StartWorkers(s)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DocEnrichTopic,
			s,
			server.SubscriberQueue("docenrich"),
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
