package main

import (
	config_handler "github.com/kazoup/platform/config/srv/handler"
	config_proto "github.com/kazoup/platform/config/srv/proto/config"
	crawler_handler "github.com/kazoup/platform/crawler/srv/handler"
	crawler_proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	crawler_subscriber "github.com/kazoup/platform/crawler/srv/subscriber"
	datasource_handler "github.com/kazoup/platform/datasource/srv/handler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_engine "github.com/kazoup/platform/db/srv/engine"
	"log"
	//_ "github.com/kazoup/platform/db/srv/engine/bleve"
	_ "github.com/kazoup/platform/db/srv/engine/elastic"
	db_handler "github.com/kazoup/platform/db/srv/handler"
	flag_handler "github.com/kazoup/platform/flag/srv/handler"
	flag_proto "github.com/kazoup/platform/flag/srv/proto/flag"
	search_engine "github.com/kazoup/platform/search/srv/engine"
	_ "github.com/kazoup/platform/search/srv/engine/db_search"
	search_handler "github.com/kazoup/platform/search/srv/handler"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"

	"github.com/kazoup/platform/structs/categories"
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

const (
	ScanTopic     = "go.micro.topic.scan"
	FileTopic     = "go.micro.topic.files"
	DbServiceName = "go.micro.srv.desktop"
)

func main() {
	cmd.Init()

	// Desktop service
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
			//parts := strings.Split(c.String("elasticsearch_hosts"), ",")
			//elastic.Hosts = parts
		}),
	)

	// Init srv
	service.Init()

	config_proto.RegisterConfigHandler(service.Server(), &config_handler.Config{
		Client:        service.Client(),
		DbServiceName: DbServiceName,
	})

	// Register DB handler Handler
	service.Server().Handle(
		service.Server().NewHandler(new(db_handler.DB)),
	)

	// Attach DB indexer subscriber
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(FileTopic, db_engine.SubscribeFiles)); err != nil {
		log.Fatal(err)
	}

	// Init DB engine
	if err := db_engine.Init(); err != nil {
		log.Fatal(err)
	}

	// Flag handler
	flag_proto.RegisterFlagHandler(service.Server(), &flag_handler.Flag{
		Client:        service.Client(),
		DbServiceName: DbServiceName,
	})

	// DataSource handler
	datasource_proto.RegisterDataSourceHandler(service.Server(), &datasource_handler.DataSource{
		Client:             service.Client(),
		ElasticServiceName: DbServiceName,
	})

	// Monitoring services
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

	if err := categories.SetMap(); err != nil {
		log.Fatal(err)
	}

	// Attach search handler
	search_proto.RegisterSearchHandler(service.Server(), &search_handler.Search{
		Client:             service.Client(),
		ElasticServiceName: DbServiceName,
	})

	if err := search_engine.Init(); err != nil {
		log.Fatalf("%s", err)
	}

	// Attach crawler handler
	crawler_proto.RegisterCrawlHandler(service.Server(), new(crawler_handler.Crawl))

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
