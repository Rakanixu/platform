package desktop

import (
	"net/http"

	"github.com/micro/cli"
	"github.com/micro/go-web"
)

func ui(ctx *cli.Context) {
	service := web.NewService(
		web.Name("com.kazoup.web.ui"),
		web.Address("0.0.0.0:80"),
	)
	service.Handle("/", http.FileServer(http.Dir("../../ui/web/html")))

	service.Run()
}

func uiCommands() []cli.Command {
	return []cli.Command{{
		Name:   "ui",
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
