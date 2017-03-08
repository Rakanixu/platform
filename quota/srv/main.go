package main

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/quota/srv/handler"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	// New service
	service := wrappers.NewKazoupService("quota", globals.QUOTA_HANDLER_QUOTA, globals.QUOTA_SUBS_QUOTA, m)

	// Monitor for file-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Healthchecks for quota-srv
	healthchecks.RegisterQuotaHealthChecks(service, m)
	healthchecks.RegisterBrokerHealthChecks(service, m)

	// New service handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Quota{
			Client: service.Client(),
		}),
	)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}