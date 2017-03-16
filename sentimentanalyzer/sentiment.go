package sentimentanalyzer

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/sentimentanalyzer/srv/handler"
	"github.com/kazoup/platform/sentimentanalyzer/srv/subscriber"
	"github.com/micro/cli"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func srv(ctx *cli.Context) {
	var m monitor.Monitor

	service := wrappers.NewKazoupService("sentimentanalyzer", m)

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
			&handler.SentimentAnalyzer{
				Client: service.Client(),
			},
		),
	); err != nil {
		log.Println(err)
	}

	s := &subscriber.SentimentAnalyzer{
		Client:        service.Client(),
		EnrichMsgChan: make(chan subscriber.EnrichMsgChan, 1000000),
		Workers:       40,
	}
	subscriber.StartWorkers(s)

	// Attach subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.SentimentEnrichTopic,
			s,
			server.SubscriberQueue("sentimentanalyzer"),
		),
	); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func sentimentAnalyzerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run sentiment analyzer service",
			Action: srv,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "sentimentanalyzer",
			Usage:       "sentiment Analyzer commands",
			Subcommands: sentimentAnalyzerCommands(),
		},
	}
}
