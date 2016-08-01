package main

import (
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/kazoup/platform/desktop/web/handler"
	"github.com/micro/go-web"
)

func main() {
	service := web.NewService(web.Name("go.micro.web.desktop"))
	service.Handle("/", http.FileServer(http.Dir("app")))
	service.Handle("/status", websocket.Handler(handler.Status))
	service.Init()
	service.Run()
}
