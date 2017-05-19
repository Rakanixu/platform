package main

import (
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/healthchecks"
	"github.com/kazoup/platform/lib/objectstorage"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"log"
	"time"
)

func main() {
	// Init DB operations
	if err := operations.Init(); err != nil {
		log.Fatal(err)
	}

	// Init custom DB operations
	if err := custom.Init(); err != nil {
		log.Fatal(err)
	}

	// Init Object Storage
	if err := objectstorage.Init(); err != nil {
		log.Fatal(err)
	}

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
