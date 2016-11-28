package handler

import (
	"fmt"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"log"
)

type Notification struct {
	Server server.Server
}

func (n *Notification) Stream(ctx context.Context, stream server.Streamer) error {
	// Listen for StreamRequest (this is blocking)
	req := &proto.StreamRequest{}
	if err := stream.Recv(req); err != nil {
		fmt.Println("ERROR receiving stream request", err)
		return err
	}

	// StreamNotifications subscribes to NotificationTopic and return channels for communications
	ch, exit, err := StreamNotifications(n.Server, req)
	if err != nil {
		fmt.Println("ERROR StreamNotifications", err)
		return err
	}

	fmt.Println("Defer stream.close", err)
	defer func() {
		close(exit)
		fmt.Println("stream.closed", err)
		stream.Close()
	}()

	for {
		select {
		// Listen over the open channel, all received notification will be pushed over this channel
		// Once channel retrieves data, send it back over the stream
		case e := <-ch:
			if err := stream.Send(&proto.StreamResponse{Message: e}); err != nil {
				log.Println("ERROR sending notification message over stream: ", err)
				return err
			}
		}
	}

	fmt.Println("Stream before return nil")

	return nil
}
