package main

import (
	"log"

	"github.com/kazoup/platform/auth/web/handler"
	web "github.com/micro/go-web"
)

func main() {

	service := web.NewService(web.Name("com.kazoup.web.auth"))
	service.HandleFunc("/google/login", handler.HandleGoogleLogin)
	service.HandleFunc("/GoogleCallback", handler.HandleGoogleCallback)
	service.HandleFunc("/google/callback", handler.HandleGoogleCallback)
	service.HandleFunc("/microsoft/login", handler.HandleMicrosoftLogin)
	service.HandleFunc("/microsoft/callback", handler.HandleMicrosoftCallback)
	service.HandleFunc("/slack/login", handler.HandleSlackLogin)
	service.HandleFunc("/slack/callback", handler.HandleSlackCallback)

	if err := service.Init(); err != nil {
		log.Panic(err)
	}

	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
