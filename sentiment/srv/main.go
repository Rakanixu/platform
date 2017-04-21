package main

import (
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/sentiment/srv/handler"
	"github.com/kazoup/platform/sentiment/srv/proto/sentiment"
	"github.com/kazoup/platform/sentiment/srv/subscriber"
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

	service := wrappers.NewKazoupService("sentiment", m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	proto_sentiment.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.SentimentEnrichTopic,
			subscriber.NewTaskHandler(40),
			server.SubscriberQueue("sentiment"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.SentimentHandler),
			server.SubscriberQueue("announce-sentiment"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
