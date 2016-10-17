package auth

import (
	"log"

	"github.com/kazoup/platform/auth/web/handler"
	"github.com/micro/cli"
	webmicro "github.com/micro/go-web"
)

func web(ctx *cli.Context) {

	service := webmicro.NewService(webmicro.Name("go.micro.web.auth"))
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
