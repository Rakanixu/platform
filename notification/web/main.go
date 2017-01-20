package main

import (
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	web_handler "github.com/kazoup/platform/notification/web/handler"
	"github.com/kazoup/platform/notification/web/sockets"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
	"time"
)

func main() {
	var m monitor.Monitor

	web := microweb.NewService(microweb.Name("com.kazoup.web.notification"))

	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)

	healthchecks.RegisterNotificationWebHealthChecks(m)

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))
	web.HandleFunc("/health", web_handler.HandleHealthCheck)

	if err := web.Init(); err != nil {
		log.Fatal(err)
	}

	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
