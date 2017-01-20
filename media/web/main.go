package main

import (
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"time"
)

func main() {
	var m monitor.Monitor

	service := microweb.NewService(microweb.Name("com.kazoup.web.media"))

	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterMediaWebHealthChecks(m)

	service.Handle("/preview", handler.NewImageHandler())
	service.Handle("/thumbnail", handler.NewThumbnailHandler())
	service.HandleFunc("/health", handler.HandleHealthCheck)

	service.Init()
	service.Run()
}
