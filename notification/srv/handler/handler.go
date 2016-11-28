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
	/*
		if len(req.UserId) == 0 {
			return errors.BadRequest("go.micro.srv.notification.Stream", "invalid user_id")
		}
	*/
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!")

	req := &proto.StreamRequest{}
	if err := stream.Recv(req); err != nil {
		fmt.Println("ERROR receiving stream request", err)
		return err
	}

	log.Println("StreamNotifications(n.Server, req)", req.UserId)

	ch, exit, err := StreamNotifications(n.Server, req)
	if err != nil {
		return err
	}

	defer func() {
		log.Println("EXIST CLOSE")
		close(exit)
		stream.Close()
	}()

	for {
		select {
		case e := <-ch:
			if err := stream.Send(&proto.StreamResponse{Message: e}); err != nil {
				log.Println("ERROR sending notification message over stream: ", err)
				return err
			}
		}
	}

	return nil
}
