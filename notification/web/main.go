package main

import (
	"log"

	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/plugins"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/notification/web/sockets"
	"github.com/micro/go-micro/client"
	microweb "github.com/micro/go-web"
	"golang.org/x/net/websocket"
)

func main() {
	web := microweb.NewService(microweb.Name("com.kazoup.web.notification"))

	// Attach socket stream
	web.Handle("/platform/notify", websocket.Handler(sockets.Stream))

	sockets.NotificationClient = proto.NewNotificationClient(
		globals.NOTIFICATION_SERVICE_NAME,
		client.DefaultClient,
	)

	if err := web.Init(); err != nil {
		log.Fatal(err)
	}

	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
