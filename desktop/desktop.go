package desktop

import (
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"net/http"
)

func ui(ctx *cli.Context) {
	service := web.NewService(web.Name("go.micro.web.ui"))
	service.Handle("/", http.FileServer(http.Dir("app")))

	service.Run()
}

func uiCommands() []cli.Command {
	return []cli.Command{{
		Name:   "desktop",
		Usage:  "Run desktop service",
		Action: ui,
	},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "ui",
			Usage:       "UI commands",
			Subcommands: uiCommands(),
		},
	}
}
