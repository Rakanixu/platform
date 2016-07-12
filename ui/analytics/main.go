package main

import (
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-web"
	"log"
	"net/http"
)

var serveFrom string

func setup(app *ccli.App) {
	app.Flags = append(
		app.Flags,
		ccli.StringFlag{
			Name:   "environment",
			EnvVar: "ENVIRONMENT",
			Usage:  "Web app serve environment (dev / prod)",
			Value:  "dev",
		},
	)

	before := app.Before

	app.Before = func(ctx *ccli.Context) error {
		if ctx.String("environment") == "prod" {
			serveFrom = "frontend/dist/sections/analytics"
		} else {
			serveFrom = "frontend/app/sections/analytics"
		}

		return before(ctx)
	}
}

func main() {
	app := cmd.App()
	setup(app)
	cmd.Init()

	// New service for serving static files
	service := web.NewService(
		web.Name("go.micro.web.analytics"),
		web.Version("0.0.1"),
	)

	// Service handlers
	service.Handle(
		"/",
		http.FileServer(http.Dir(serveFrom)),
	)

	// Init and run service
	if err := service.Init(); err != nil {
		log.Fatalf("%v", err)
	}

	if err := service.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
