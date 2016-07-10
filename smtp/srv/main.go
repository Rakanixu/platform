package main

import (
	"log"

	"github.com/kazoup/smtp/srv/handler"
	"github.com/kazoup/smtp/srv/smtp"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.smtp"),
		micro.Version("latest"),
		micro.Flags(
			cli.StringFlag{
				Name:   "email_host",
				EnvVar: "EMAIL_HOST",
				Usage:  "SMTP server address",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "email_host_port",
				EnvVar: "EMAIL_HOST_PORT",
				Usage:  "SMTP server port",
				Value:  "25",
			},
			cli.StringFlag{
				Name:   "email_host_user",
				EnvVar: "EMAIL_HOST_USER",
				Usage:  "User",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "email_host_password",
				EnvVar: "EMAIL_HOST_PASSWORD",
				Usage:  "Password",
				Value:  "",
			},
			cli.StringFlag{
				Name:   "default_from_email",
				EnvVar: "DEFAULT_FROM_EMAIL",
				Usage:  "Default mail address sender",
				Value:  "noreply@org.com",
			},
		),
		micro.Action(func(c *cli.Context) {
			smtp.EmailHost = c.String("email_host")
			smtp.EmailHostPort = c.String("email_host_port")
			smtp.EmailHostUser = c.String("email_host_user")
			smtp.EmailHostPassword = c.String("email_host_password")
			smtp.DefaultFromEmail = c.String("default_from_email")
		}),
	)

	// Register Handler
	service.Server().Handle(
		service.Server().NewHandler(new(handler.SMTP)),
	)

	// Initialise service
	service.Init()

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
