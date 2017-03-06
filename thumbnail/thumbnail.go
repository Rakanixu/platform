package enrich

import (
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/thumbnail/srv/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("thumbnail", globals.QUOTA_HANDLER_THUMBNAIL, globals.QUOTA_SUBS_THUMBNAIL, m)

	// enrich-srv monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	gcslib.Register()

	s := &subscriber.Thumbnail{
		Client:             service.Client(),
		GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		ThumbnailMsgChan:   make(chan *enrich_proto.EnrichMessage, 1000000),
		Workers:            25,
	}
	subscriber.StartWorkers(s)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.ThumbnailTopic,
			s,
			server.SubscriberQueue("thumbnail"),
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func thumbnailCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run thumbnail service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "thumbnail",
			Usage:       "Thumbnail commands",
			Subcommands: thumbnailCommands(),
		},
	}
}
