package main

import (
	"github.com/kazoup/platform/agent/srv/handler"
	"github.com/kazoup/platform/agent/srv/proto/agent"
	"github.com/kazoup/platform/agent/srv/subscriber"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/monitor"
	"log"
	"time"
)

// Runs the agent service
func main() {
	var m monitor.Monitor

	// Create instance of agent service
	service := wrappers.NewKazoupService("agent", m)

	// Create monitor for agent service
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
		monitor.Client(service.Client()),
		monitor.Server(service.Server()),
	)
	defer m.Close()

	// Register broker health checks
	healthchecks.RegisterBrokerHealthChecks(service, m)

	// Attach handler
	proto_agent.RegisterServiceHandler(service.Server(), new(handler.Service))

	// Attach subscriber for task handler topics
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.SaveRemoteFilesTopic,
			new(subscriber.AgentServiceTaskHandler),
			server.SubscriberQueue("agent"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Attach subscriber for announce topic
	if err := service.Server().Subscribe(
		service.Server().NewSubscriber(
			globals.AnnounceTopic,
			new(subscriber.AnnounceHandler),
			server.SubscriberQueue("agent"),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Start the server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
