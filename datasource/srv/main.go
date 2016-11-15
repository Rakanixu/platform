package main

import (
	"log"

	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	service := wrappers.NewKazoupService("datasource")

	// Init broker on subscriber
	// This is required to be able to handle the data properly when
	// we want to stream the messages over the notification socket
	subscriber.Broker = service.Server().Options().Broker

	// Attach crawler started subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerStartedTopic, subscriber.SubscribeCrawlerStarted)); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, subscriber.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.DataSource{
			Client: service.Client(),
		}),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
