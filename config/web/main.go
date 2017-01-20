package main

import (
	"github.com/kazoup/platform/config/web/handler"
	"github.com/kazoup/platform/config/web/sockets"
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	web := microweb.NewService(microweb.Name("com.kazoup.web.config"))

	// config-web monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterConfigWebHealthChecks(m)

	// Attach web handler (socket)
	web.Handle("/platform/ping", websocket.Handler(sockets.PingPlatform))
	web.HandleFunc("/health", handler.HandleHealthCheck)

	web.Init()
	// Run service
	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
