package main

import (
	"github.com/kazoup/platform/lib/db/operations"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/user/srv/handler"
	"github.com/kazoup/platform/user/srv/proto/user"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	// New service
	service := wrappers.NewKazoupService("user", m)

	// Monitor for file-srv
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// New service handler
	proto_user.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Init DB operations
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
