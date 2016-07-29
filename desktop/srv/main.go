package main

import (
	"log"
	"strings"

	crawler "github.com/kazoup/platform/crawler/srv/handler"
	crawler_proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	crawler_subscriber "github.com/kazoup/platform/crawler/srv/subscriber"
	elastic "github.com/kazoup/platform/elastic/srv/elastic"
	search "github.com/kazoup/platform/elastic/srv/handler"
	indexer "github.com/kazoup/platform/elastic/srv/subscriber"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/broker/mock"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry/mdns"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/micro/monitor-srv/handler"
	"github.com/micro/monitor-srv/monitor"
	"github.com/micro/monitor-srv/proto/monitor"
)

const ScanTopic string = "go.micro.topic.scan"
const FileTopic string = "go.micro.topic.files"

func main() {
	service := micro.NewService(
		micro.Name("go.micro.srv.desktop"),
		micro.Version("latest"),
		micro.Broker(mock.NewBroker()),
		micro.Registry(mdns.NewRegistry()),
		micro.Flags(
			cli.StringFlag{
				Name:   "elasticsearch_hosts",
				EnvVar: "ELASTICSEARCH_HOSTS",
				Usage:  "Comma separated list of elasticsearch hosts",
				Value:  "localhost:9200",
			},
		),
		micro.Action(func(c *cli.Context) {
			parts := strings.Split(c.String("elasticsearch_hosts"), ",")
			elastic.Hosts = parts
		}),
	)

	cmd.Init()
	// Monitoring service

	// healthchecks
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.HealthCheckTopic,
			monitor.DefaultMonitor.ProcessHealthCheck,
		),
	)

	// status
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatusTopic,
			monitor.DefaultMonitor.ProcessStatus,
		),
	)

	// stats
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatsTopic,
			monitor.DefaultMonitor.ProcessStats,
		),
	)

	go_micro_srv_monitor_monitor.RegisterMonitorHandler(service.Server(), new(handler.Monitor))
	// Init srv
	service.Init()
	//Init elasticsearch
	elastic.Init()
	// Register search handler
	service.Server().Handle(
		service.Server().NewHandler(new(search.Elastic)),
	)

	// Attach crawler handler
	crawler_proto.RegisterCrawlHandler(service.Server(), new(crawler.Crawl))

	// Attach indexer subsciber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, indexer.FileSubscriber),
	); err != nil {
		log.Fatal(err)
	}

	// Attach crawler subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			ScanTopic,
			crawler_subscriber.Scans,
		),
	); err != nil {
		log.Fatal(err)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
