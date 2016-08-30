package auth

import (
	"log"

	"github.com/kazoup/platform/auth/web/handler"
	webmicro "github.com/micro/go-web"
	"github.com/micro/cli"
)

func web(ctx *cli.Context) {

	service := webmicro.NewService(webmicro.Name("go.micro.web.auth"))
	service.HandleFunc("/google/login", handler.HandleGoogleLogin)
	service.HandleFunc("/GoogleCallback", handler.HandleGoogleCallback)
	service.HandleFunc("/google/callback", handler.HandleGoogleCallback)
	service.HandleFunc("/microsoft/login", handler.HandleMicrosoftLogin)
	service.HandleFunc("/microsoft/callback", handler.HandleMicrosoftCallback)
	service.HandleFunc("/slack/login", handler.HandleSlackLogin)
	service.HandleFunc("/slack/callback", handler.HandleSlackCallback)

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
