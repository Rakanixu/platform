package main

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("datasource", m)

	// datasource-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterDatasourceHealthChecks(service, m)
	healthchecks.RegisterBrokerHealthChecks(service, m)

	gcslib.Register()

	// Attach crawler started subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.AnnounceTopic, &subscriber.Announce{
			Client: service.Client(),
			Broker: service.Server().Options().Broker,
		})); err != nil {
		log.Fatal(err)
	}

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
			Client:             service.Client(),
			Broker:             service.Server().Options().Broker,
			GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		})); err != nil {
		log.Fatal(err)
	}

	// Attach clean bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(globals.DeleteFileInBucketTopic, &subscriber.DeleteFileInBucket{
			Client:             service.Client(),
			Broker:             service.Server().Options().Broker,
			GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		})); err != nil {
		log.Fatal(err)
	}

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.DataSource{
			Client:             service.Client(),
			GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		}),
	)
	service.Init()
	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
