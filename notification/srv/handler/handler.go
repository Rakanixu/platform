package handler

import (
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

type Notification struct {
	Server server.Server
}

func (n *Notification) Stream(ctx context.Context, req *proto.StreamRequest, stream proto.Notification_StreamStream) error {
	if len(req.UserId) == 0 {
		return errors.BadRequest("go.micro.srv.message", "invalid user_id")
	}

	ch, exit, err := StreamNotifications(n.Server, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.message", err.Error())
	}

	defer func() {
		close(exit)
		stream.Close()
	}()

	for {
		select {
		case e := <-ch:
			if err := stream.Send(&proto.StreamResponse{Message: e}); err != nil {
				return err
			}
		}
	}

	return nil
}
