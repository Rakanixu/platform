package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/kardianos/osext"
	auth "github.com/kazoup/platform/auth"
	config "github.com/kazoup/platform/config"
	crawler "github.com/kazoup/platform/crawler"
	datasource "github.com/kazoup/platform/datasource"
	db "github.com/kazoup/platform/db"
	file "github.com/kazoup/platform/file"
	"github.com/kazoup/platform/lib/globals"
	media "github.com/kazoup/platform/media"
	monitor "github.com/kazoup/platform/monitor"
	notification "github.com/kazoup/platform/notification"
	search "github.com/kazoup/platform/search"
	ui "github.com/kazoup/platform/ui"
	"github.com/micro/cli"
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/micro/web"
)

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, auth.Commands()...)
	app.Commands = append(app.Commands, config.Commands()...)
	app.Commands = append(app.Commands, crawler.Commands()...)
	app.Commands = append(app.Commands, datasource.Commands()...)
	app.Commands = append(app.Commands, db.Commands()...)
	app.Commands = append(app.Commands, media.Commands()...)
	app.Commands = append(app.Commands, search.Commands()...)
	app.Commands = append(app.Commands, file.Commands()...)
	app.Commands = append(app.Commands, notification.Commands()...)
	app.Commands = append(app.Commands, monitor.Commands()...)
	app.Commands = append(app.Commands, ui.Commands()...)
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
		ccli.BoolFlag{
			Name:   "enable_tls",
			Usage:  "Enable TLS",
			EnvVar: "MICRO_ENABLE_TLS",
		},
		ccli.StringFlag{
			Name:   "tls_cert_file",
			Usage:  "TLS Certificate file",
			EnvVar: "MICRO_TLS_CERT_FILE",
		},
		ccli.StringFlag{
			Name:   "tls_key_file",
			Usage:  "TLS Key file",
			EnvVar: "MICRO_TLS_KEY_FILE",
		},
		ccli.StringFlag{
			Name:   "tls_client_ca_file",
			Usage:  "TLS CA file to verify clients against",
			EnvVar: "MICRO_TLS_CLIENT_CA_FILE",
		},
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
	var nc *exec.Cmd
	cmds := ctx.App.Commands
	binary, err := osext.Executable()
	if err != nil {
		log.Println(err.Error())
	}
	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Println(err.Error())
	}

	// Execute nats server binary.
	wg.Add(1)
	nc = exec.Command(
		fmt.Sprintf("%s%s%s%s%s/gnatsd", dir, "/nats/gnatsd-v0.9.4-", runtime.GOOS, "-", runtime.GOARCH),
		"-m",
		"8222",
	)
	nc.Stdout = os.Stdout
	nc.Stderr = os.Stderr
	if err := nc.Start(); err != nil {
		log.Println(err.Error())
		wg.Done()
	}
	time.Sleep(time.Second * 2)

	// Execute consul as registry in development
	// consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul -bind=127.0.0.1
	/*
		var cr *exec.Cmd
		wg.Add(1)
		cr = exec.Command(
			fmt.Sprintf("%s%s%s%s%s/consul", dir, "/consul/consul_0.7.2_", runtime.GOOS, "_", runtime.GOARCH),
			"agent",
			"-server",
			"-bootstrap-expect",
			"1",
			"-data-dir",
			"/tmp/consul",
			"-bind=127.0.0.1",
		)
		//cr.Stdout = os.Stdout
		//cr.Stderr = os.Stderr
		if err := cr.Start(); err != nil {
			log.Fatal(err.Error())
			wg.Done()
		}
		time.Sleep(time.Second * 3)
	*/

	for _, cmd := range cmds {
		if cmd.Name != "help" && len(cmd.Subcommands) > 0 {
			for _, subcmd := range cmd.Subcommands {
				wg.Add(1)
				c := exec.Command(
					binary,
					"--registry=consul",
					"--registry_address=127.0.0.1",
					"--transport=tcp",
					"--broker=nats",
					"--broker_address=127.0.0.1:4222",
					"--enable_tls",
					"--tls_cert_file=ssl/all.pem",
					"--tls_key_file=ssl/key.pem",
					cmd.Name,
					subcmd.Name,
				)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if err := c.Start(); err != nil {
					log.Println(err.Error())
					wg.Done()
				}
			}
		}
		if cmd.Name != "help" && len(cmd.Subcommands) == 0 && cmd.Name != "desktop" {
			wg.Add(1)
			c := exec.Command(
				binary,
				"--registry=consul",
				"--registry_address=127.0.0.1",
				"--broker=nats",
				"--transport=tcp",
				"--broker_address=127.0.0.1:4222",
				"--enable_tls",
				"--tls_cert_file=ssl/all.pem",
				"--tls_key_file=ssl/key.pem",
				cmd.Name)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Start(); err != nil {
				log.Println(err.Error())
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
