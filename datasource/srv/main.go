package main

import (
	"github.com/kazoup/platform/datasource/srv/handler"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/datasource/srv/subscriber"
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
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("announce-datasource"),
		)); err != nil {
		log.Fatal(err)
	}

	// Attach crawler finished subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DiscoveryFinishedTopic,
			new(subscriber.DiscoveryFinished),
			server.SubscriberQueue("discoveryfinished"),
		)); err != nil {
		log.Fatal(err)
	}

	// Attach delete bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DeleteBucketTopic,
			subscriber.NewDeleteBucketHandler(gcslib.NewGoogleCloudStorage()),
			server.SubscriberQueue("deletebucket"),
		)); err != nil {
		log.Fatal(err)
	}

	// Attach clean bucket subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.DeleteFileInBucketTopic,
			subscriber.NewDeleteFileInBucketHandler(gcslib.NewGoogleCloudStorage()),
			server.SubscriberQueue("deletefileinbucket"),
		)); err != nil {
		log.Fatal(err)
	}

	// New service handler
	proto_datasource.RegisterServiceHandler(service.Server(), handler.NewServiceHandler(gcslib.NewGoogleCloudStorage()))

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
