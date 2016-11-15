package main

import (
	"log"

	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/notification/web/sockets"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/tcp"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
)

func main() {
	web := microweb.NewService(microweb.Name("com.kazoup.web.notification"))

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))

	sockets.NotificationClient = proto.NewNotificationClient(
		"com.kazoup.srv.notification",
		client.DefaultClient,
	)

	if err := web.Init(); err != nil {
		log.Fatal(err)
	}

	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
