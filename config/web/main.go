package main

import (
	"golang.org/x/net/websocket"
	microweb "github.com/micro/go-web"
	"github.com/kazoup/platform/config/srv/sockets"
	"log"
	_ "github.com/micro/go-plugins/broker/nats"
)

func main() {
	web := microweb.NewService(microweb.Name("com.kazoup.web.config"))

	// Attach web handler (socket)
	web.Handle("/platform/ping", websocket.Handler(sockets.PingPlatform))


    web.Init()
	// Run service
	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
