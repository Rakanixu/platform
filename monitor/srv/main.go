package main

import (
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/monitor/srv/handler"
	"github.com/kazoup/platform/monitor/srv/monitor"
	proto "github.com/kazoup/platform/monitor/srv/proto/monitor"
	"github.com/micro/go-micro"
	os_monitor "github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m os_monitor.Monitor

	service := micro.NewService(
		micro.Name(globals.MONITOR_SERVICE_NAME),
		micro.BeforeStart(func() error {
			monitor.DefaultMonitor.Run()
			return nil
		}),
	)

	m = os_monitor.NewMonitor(
		os_monitor.Interval(time.Minute),
		os_monitor.Client(service.Client()),
		os_monitor.Server(service.Server()),
	)
	defer m.Close()

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

	proto.RegisterMonitorHandler(service.Server(), new(handler.Monitor))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
