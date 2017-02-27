package enrich

import (
	"github.com/kazoup/platform/audioenrich/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/cli"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
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

	gcslib.Register()

	s := &subscriber.Enrich{
		Client:             service.Client(),
		GoogleCloudStorage: gcslib.NewGoogleCloudStorage(),
		EnrichMsgChan:      make(chan *enrich_proto.EnrichMessage, 100000),
	}
	subscriber.SyncMessages(s)

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

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func audioEnrichCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run audio enrich service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "audioenrich",
			Usage:       "Audio Enrich commands",
			Subcommands: audioEnrichCommands(),
		},
	}
}
