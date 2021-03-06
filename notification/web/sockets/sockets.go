package sockets

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/websocket"
)

func Stream(ws *websocket.Conn) {
	var m map[string]interface{}
	c := client.NewClient()

	// Connection established by client, server responds to let client know it can start sending data
	if err := websocket.JSON.Send(ws, struct {
		Connected bool `json:"connected"`
	}{Connected: true}); err != nil {
		fmt.Println("ERROR sending connected state", err)
		return
	}

	// Client sends UserID
	if err := websocket.JSON.Receive(ws, &m); err != nil {
		fmt.Println("ERROR receiving user_id /notificaion/platform/notify", err)
		return
	}

	// Stream initialization
	sreq := client.DefaultClient.NewRequest(
		globals.NOTIFICATION_SERVICE_NAME,
		"Service.Stream",
		&proto_notification.StreamRequest{},
	)

	if m["token"] == nil {
		fmt.Println("ERROR receiving jwt_token from client", "jwt_token not found")
		return
	}

	stream, err := c.Stream(globals.NewContextFromJWT(m["token"].(string)), sreq)
	if err != nil {
		fmt.Println("ERROR opening stream for notifications", err)
		return
	}

	defer stream.Close()

	// Send to Notification srv the userID we received from client connection
	// At this moment,we subscribe to NotificationTopic
	// Once is subscribed, service does not expect to received more data from client
	if err := stream.Send(&proto_notification.StreamRequest{
		UserId: m["user_id"].(string),
	}); err != nil {
		fmt.Println("ERROR sending userID to stream handler")
		return
	}

	// Listen for StreamResponses from notification service
	// Once a response is received, send it back to client over the socket connection
	for {
		srsp := &proto_notification.StreamResponse{}
		if err := stream.Recv(srsp); err != nil {
			fmt.Println("ERROR receiving notification from stream", err)
			return
		}

		if err := websocket.JSON.Send(ws, srsp.Message); err != nil {
			fmt.Println("ERROR sending notification over websocket", err)
			return
		}

	}
}
