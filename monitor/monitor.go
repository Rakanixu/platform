package monitor

import (
	"fmt"
	"github.com/kardianos/osext"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/monitor/srv/handler"
	"github.com/kazoup/platform/monitor/srv/monitor"
	proto "github.com/kazoup/platform/monitor/srv/proto/monitor"
	web_handler "github.com/kazoup/platform/monitor/web/handler"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	microweb "github.com/micro/go-web"
	"log"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name(globals.MONITOR_SERVICE_NAME),
		// before starting
		micro.BeforeStart(func() error {
			monitor.DefaultMonitor.Run()
			return nil
		}),
	)

	// healthchecks
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.HealthCheckTopic,
			monitor.DefaultMonitor.ProcessHealthCheck,
		),
	)

	// status
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatusTopic,
			monitor.DefaultMonitor.ProcessStatus,
		),
	)

	// stats
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatsTopic,
			monitor.DefaultMonitor.ProcessStats,
		),
	)

	proto.RegisterMonitorHandler(service.Server(), new(handler.Monitor))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Println(err.Error())
	}

	service := microweb.NewService(
		microweb.Name("go.micro.web.monitor"),
		microweb.Handler(web_handler.Router()),
	)

	web_handler.Init(
		fmt.Sprintf("%s%s", dir, "../../../monitor/web/templates"),
		proto.NewMonitorClient(globals.MONITOR_SERVICE_NAME, client.DefaultClient),
	)

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}

func configCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run monitor srv service",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run monitor web service",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "monitor",
			Usage:       "Monitor commands",
			Subcommands: configCommands(),
		},
	}
}
