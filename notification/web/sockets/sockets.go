package sockets

import (
	"fmt"

	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/websocket"
)

//var NotificationClient proto.NotificationClient

func Stream(ws *websocket.Conn) {
	var m map[string]interface{}

	if err := websocket.JSON.Send(ws, struct {
		Connected bool `json:"connected"`
	}{Connected: true}); err != nil {
		fmt.Println("ERROR sending connected state", err)
		return
	}

	if err := websocket.JSON.Receive(ws, &m); err != nil {
		fmt.Println("ERROR receiving user_id /notificaion/platform/notify", err)
		return
	}

	fmt.Println("MSG received", m["user_id"].(string))

	sreq := client.DefaultClient.NewRequest(
		globals.NOTIFICATION_SERVICE_NAME,
		"Notification.Stream",
		&proto.StreamRequest{
		/*UserId: m["user_id"].(string),*/
		},
	)

	stream, err := client.DefaultClient.Stream(globals.NewSystemContext(), sreq)
	if err != nil {
		fmt.Println("ERROR opening stream for notifications", err)
		return
	}

	defer stream.Close()

	if err := stream.Send(&proto.StreamRequest{
		UserId: m["user_id"].(string),
	}); err != nil {
		fmt.Println("ERROR sending userID to stream handler")
		return
	}

	for {
		srsp := &proto.StreamResponse{}
		/*msg, */ err := stream.Recv(srsp)
		fmt.Println("YEII", srsp)

		if err != nil {
			fmt.Println("ERROR receiving notification from stream", err)
			return
		}

		if err := websocket.JSON.Send(ws, srsp.Message); err != nil {
			fmt.Println("ERROR sending notification over websocket", err)
			return
		}

	}
}
