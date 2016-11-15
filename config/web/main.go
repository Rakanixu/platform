package main

import (
	"log"

	"github.com/kazoup/platform/config/srv/sockets"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
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
