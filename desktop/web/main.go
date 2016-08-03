package main

import (
	"net/http"

	"golang.org/x/net/websocket"

	"github.com/kazoup/platform/desktop/web/handler"
	"github.com/micro/go-web"
)

func main() {
	contentDir := "/"
	service := web.NewService(web.Name("go.micro.web.desktop"))
	service.Handle("/", http.FileServer(http.Dir("app")))
	service.Handle("/status", websocket.Handler(handler.Status))
	service.Handle("/stream/", http.StripPrefix("/stream/", handler.NewPlaylistHandler(contentDir)))
	service.Handle("/frame/", http.StripPrefix("/frame/", handler.NewFrameHandler(contentDir)))
	service.Handle("/segments/", http.StripPrefix("/segments/", handler.NewStreamHandler(contentDir)))

	service.Handle("/webm/", http.StripPrefix("/webm/", handler.NewWebmHandler(contentDir)))
	service.Init()
	service.Run()
}
