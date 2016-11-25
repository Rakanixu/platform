package handler

import (
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type Notification struct {
	Service micro.Service
}

func (n *Notification) Stream(ctx context.Context, req *proto.StreamRequest, stream proto.Notification_StreamStream) error {
	if len(req.UserId) == 0 {
		return errors.BadRequest("go.micro.srv.notification.Stream", "invalid user_id")
	}

	ch, exit, err := StreamNotifications(n.Service, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.notification.StreamNotifications: ", err.Error())
	}

	defer func() {
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
