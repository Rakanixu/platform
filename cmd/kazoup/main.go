package main

import (
	"github.com/kardianos/osext"
	auth "github.com/kazoup/platform/auth"
	config "github.com/kazoup/platform/config"
	crawler "github.com/kazoup/platform/crawler"
	datasource "github.com/kazoup/platform/datasource"
	db "github.com/kazoup/platform/db"
	ui "github.com/kazoup/platform/desktop"
	flag "github.com/kazoup/platform/flag"
	media "github.com/kazoup/platform/media"
	//notification "github.com/kazoup/platform/notification"
	scheduler "github.com/kazoup/platform/scheduler"
	search "github.com/kazoup/platform/search"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/cli"
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/web"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, auth.Commands()...)
	app.Commands = append(app.Commands, config.Commands()...)
	app.Commands = append(app.Commands, crawler.Commands()...)
	app.Commands = append(app.Commands, datasource.Commands()...)
	app.Commands = append(app.Commands, db.Commands()...)
	app.Commands = append(app.Commands, ui.Commands()...)
	app.Commands = append(app.Commands, flag.Commands()...)
	app.Commands = append(app.Commands, media.Commands()...)
	app.Commands = append(app.Commands, search.Commands()...)
	app.Commands = append(app.Commands, scheduler.Commands()...)
	//app.Commands = append(app.Commands, notification.Commands()...)
	app.Commands = append(app.Commands, web.Commands()...)
	app.Commands = append(app.Commands, desktopCommands()...)
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
		ccli.StringFlag{
			Name:  "web_namespace",
			Value: globals.NAMESPACE,
		},
	)
}

func desktop(ctx *ccli.Context) {
	var wg sync.WaitGroup
	cmds := ctx.App.Commands
	binary, _ := osext.Executable()
	for _, cmd := range cmds {
		if cmd.Name != "help" && len(cmd.Subcommands) > 0 {
			for _, subcmd := range cmd.Subcommands {
				//time.Sleep(time.Second)
				wg.Add(1)
				log.Print(cmd.Name, subcmd.Name)
				c := exec.Command(binary, "--registry=mdns", cmd.Name, subcmd.Name)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Start(); err != nil {
					log.Print(err.Error())
					wg.Done()
				}
			}
		}
		if cmd.Name != "help" && len(cmd.Subcommands) == 0 && cmd.Name != "desktop" {

			wg.Add(1)
			log.Print(cmd.Name)
			c := exec.Command(binary, "--registry=mdns", cmd.Name)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Start(); err != nil {
				log.Print(err.Error())
				wg.Done()
			}
		}
	}
	//TODO:Start elasticsearch

	wg.Wait()

}

func desktopCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "desktop",
			Usage:  "Run desktop service",
			Action: desktop,
		},
	}
}
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "desktop",
			Usage:       "Auth commands",
			Subcommands: desktopCommands(),
		},
	}
}
