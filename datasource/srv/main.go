package main

import (
	"log"

	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
)

func main() {
	service := wrappers.NewKazoupService("datasource")

	// Attach crawler started subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerStartedTopic, &subscriber.CrawlerStarted{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.CrawlerFinishedTopic, &subscriber.CrawlerFinished{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach delete bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteBucketTopic, &subscriber.DeleteBucket{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

	// Attach clean bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteFileInBucketTopic, &subscriber.DeleteFileInBucket{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
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
