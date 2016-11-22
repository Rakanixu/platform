package sockets

import (
	"fmt"

	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"golang.org/x/net/websocket"
)

var NotificationClient proto.NotificationClient

func Stream(ws *websocket.Conn) {
	var m map[string]interface{}

	if err := websocket.JSON.Receive(ws, &m); err != nil {
		fmt.Println("err", err)
		return
	}

	stream, err := NotificationClient.Stream(globals.NewSystemContext(), &proto.StreamRequest{
		UserId: m["user_id"].(string),
	})
	if err != nil {
		fmt.Println("err", err)
		return
	}

	defer stream.Close()

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Println("err", err)
			return
		}

		if err := websocket.JSON.Send(ws, msg.Message); err != nil {
			fmt.Println("err", err)
			return
		}

	}
}