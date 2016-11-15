package main

import (
	"log"

	"github.com/kazoup/platform/lib/wrappers"
	"github.com/kazoup/platform/scheduler/srv/handler"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	// New Service
	service := wrappers.NewKazoupService("scheduler")

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(&handler.Scheduler{
			Client: service.Client(),
			Crons:  make([]*handler.CronWrapper, 0),
		}),
	)
	// Init service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
