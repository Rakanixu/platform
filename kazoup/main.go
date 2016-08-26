package main

import (
	auth "github.com/kazoup/platform/auth"
	config "github.com/kazoup/platform/config"
	crawler "github.com/kazoup/platform/crawler"
	datasource "github.com/kazoup/platform/datasource"
	elastic "github.com/kazoup/platform/elastic"
	"github.com/micro/cli"
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
)

func main() {
	app := cmd.App()

	app.Commands = append(app.Commands, auth.Commands()...)
	app.Commands = append(app.Commands, config.Commands()...)
	app.Commands = append(app.Commands, crawler.Commands()...)
	app.Commands = append(app.Commands, datasource.Commands()...)
	app.Commands = append(app.Commands, elastic.Commands()...)
	app.Action = func(context *cli.Context) { cli.ShowAppHelp(context) }

	setup(app)
	cmd.Init(
		cmd.Name("kazoup"),
		cmd.Description("Kazoup platform"),
		cmd.Version("latest"),
	)
}

func setup(app *ccli.App) {
	// common flags
	app.Flags = append(app.Flags,
		ccli.IntFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		ccli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
		},
		ccli.StringFlag{
			Name:   "html_dir",
			EnvVar: "MICRO_HTML_DIR",
			Usage:  "The html directory for a web app",
		},
		ccli.StringFlag{
			Name:   "elasticsearch_hosts",
			EnvVar: "ELASTICSEARCH_HOSTS",
			Usage:  "ELasticsearch hosts ie: localhost:9200",
		},
	)
}
