package main

import (
	"log"

	"github.com/kazoup/platform/auth/web/handler"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	web "github.com/micro/go-web"
)

func main() {

	service := web.NewService(web.Name("com.kazoup.web.auth"))

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

	if err := service.Init(); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}