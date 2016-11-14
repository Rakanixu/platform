package main

import (
	"github.com/kazoup/platform/notification/web/sockets"
	_ "github.com/micro/go-plugins/broker/nats"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	web := microweb.NewService(microweb.Name("com.kazoup.web.notification"))

	// Attach web handler (socket)
	web.Handle("/platform/ping", websocket.Handler(sockets.Stream))

	web.Init()
	// Run service
	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
