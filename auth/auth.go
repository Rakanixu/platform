package auth

import (
	"github.com/kazoup/platform/auth/web/handler"
	"github.com/kazoup/platform/lib/healthchecks"
	"github.com/micro/cli"
	"github.com/micro/go-os/monitor"
	webmicro "github.com/micro/go-web"
	"log"
	"time"
)

func web(ctx *cli.Context) {
	var m monitor.Monitor

	service := webmicro.NewService(webmicro.Name("go.micro.web.auth"))

	// auth-web monitor
	m = monitor.NewMonitor(
		monitor.Interval(time.Minute),
	)
	defer m.Close()

	healthchecks.RegisterAuthWebHealthChecks(m)

	service.HandleFunc("/google/login", handler.HandleGoogleLogin)
	service.HandleFunc("/google/callback", handler.HandleGoogleCallback)
	service.HandleFunc("/microsoft/login", handler.HandleMicrosoftLogin)
	service.HandleFunc("/microsoft/callback", handler.HandleMicrosoftCallback)
	service.HandleFunc("/slack/login", handler.HandleSlackLogin)
	service.HandleFunc("/slack/callback", handler.HandleSlackCallback)
	service.HandleFunc("/dropbox/login", handler.HandleDropboxLogin)
	service.HandleFunc("/dropbox/callback", handler.HandleDropboxCallback)
	service.HandleFunc("/box/login", handler.HandleBoxLogin)
	service.HandleFunc("/box/callback", handler.HandleBoxCallback)
	service.HandleFunc("/gmail/login", handler.HandleGmailLogin)
	service.HandleFunc("/gmail/callback", handler.HandleGmailCallback)
	service.HandleFunc("/health", handler.HandleHealthCheck)

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}

func authCommands() []cli.Command {
	return []cli.Command{{
		Name:   "web",
		Usage:  "Run auth web service",
		Action: web,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "auth",
			Usage:       "Auth commands",
			Subcommands: authCommands(),
		},
	}
}
