package main

import (
	_ "github.com/kazoup/platform/lib/plugins"
	"github.com/kazoup/platform/notification/web/sockets"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	web := microweb.NewService(microweb.Name("com.kazoup.web.notification"))

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))

	if err := web.Init(); err != nil {
		log.Fatal(err)
	}

	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
