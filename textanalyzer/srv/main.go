package main

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/textanalyzer/srv/handler"
	"github.com/kazoup/platform/textanalyzer/srv/subscriber"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("textanalyzer", globals.QUOTA_HANDLER_TEXT_ANALYZER, globals.QUOTA_SUBS_TEXT_ANALYZER, m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	if err := service.Server().Handle(
		service.Server().NewHandler(
			&handler.TextAnalyzer{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Println(err)
	}

	s := &subscriber.TextAnalyzer{
		Client:        service.Client(),
		EnrichMsgChan: make(chan subscriber.EnrichMsgChan, 1000000),
		Workers:       40,
	}
	subscriber.StartWorkers(s)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.ExtractEntitiesTopic,
			s,
			server.SubscriberQueue("textanalyzer"),
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
