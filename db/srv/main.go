package main

import (
	"github.com/kazoup/platform/db/srv/engine"
	//_ "github.com/kazoup/platform/db/srv/engine/bleve"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	"github.com/kazoup/platform/db/srv/handler"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	"log"
)

func main() {
	// New Service
	service := wrappers.NewKazoupService("db")

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.DB)),
	)

	// Attach indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.FilesTopic, engine.SubscribeFiles)); err != nil {
		log.Fatal(err)
	}

	// Attach slack user indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.SlackUsersTopic, engine.SubscribeSlackUsers)); err != nil {
		log.Fatal(err)
	}

	// Attach slack channel indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.SlackChannelsTopic, engine.SubscribeSlackChannels)); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, engine.SubscribeCrawlerFinished)); err != nil {
		log.Fatal(err)
	}

	// Initialise service
	service.Init()
	// Init search engine

	if err := engine.Init(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
