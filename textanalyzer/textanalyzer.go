package enrich

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/textanalyzer/srv/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
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

	s := &subscriber.TextAnalyzer{
		Client:        service.Client(),
		EnrichMsgChan: make(chan *enrich_proto.EnrichMessage, 1000000),
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

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func textAnalyzerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run text analyzer service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "textanalyzer",
			Usage:       "Text Analyzer commands",
			Subcommands: textAnalyzerCommands(),
		},
	}
}
