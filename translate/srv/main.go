package main

import (
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/translate/srv/handler"
	"github.com/kazoup/platform/translate/srv/proto/translate"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor
	service := wrappers.NewKazoupService("translate", m)

	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	healthchecks.RegisterBrokerHealthChecks(service, m)

	proto_translate.RegisterServiceHandler(service.Server(), new(handler.Service))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
