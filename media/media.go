package media

import (
	"github.com/kazoup/platform/lib/healthchecks"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/media/web/handler"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	microweb "github.com/micro/go-web"
	"time"
)

func web(ctx *cli.Context) {
	var m monitor.Monitor

	service := microweb.NewService(microweb.Name("go.micro.web.media"))

	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterMediaWebHealthChecks(m)

	service.Handle("/preview", handler.NewImageHandler())
	service.Handle("/thumbnail", handler.NewThumbnailHandler())
	service.HandleFunc("/health", handler.HandleHealthCheck)

	service.Run()
}

func mediaCommands() []cli.Command {
	return []cli.Command{{
		Name:   "web",
		Usage:  "Run media web service",
		Action: web,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "media",
			Usage:       "Media commands",
			Subcommands: mediaCommands(),
		},
	}
}
