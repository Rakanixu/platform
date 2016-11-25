package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"log"
)

type Notification struct {
	Server server.Server
}

func (n *Notification) Stream(ctx context.Context, req *proto.StreamRequest, stream proto.Notification_StreamStream) error {
	if len(req.UserId) == 0 {
		return errors.BadRequest("go.micro.srv.notification.Stream", "invalid user_id")
	}

	_, err := n.Server.Options().Broker.Subscribe(globals.NotificationTopic, func(p broker.Publication) error {
		var e *proto.NotificationMessage

		if err := json.Unmarshal(p.Message().Body, &e); err != nil {
			return err
		}

		if req.UserId == e.UserId {
			//che <- e

			if err := stream.Send(&proto.StreamResponse{Message: e}); err != nil {
				log.Println("ERROR sending notification message over stream: ", err)
				return err
			}
		}

		return nil
	})

	if err := stream.Send(&proto.StreamResponse{
		Message: &proto.NotificationMessage{
			Info: "INFO1",
		},
	}); err != nil {
		log.Println("INFO1", err)
	}

	//ch, exit, err := StreamNotifications(n.Server, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.notification.StreamNotifications: ", err.Error())
	}

	/*	defer func() {
		close(exit)
		stream.Close()
	}()*/

	if err := stream.Send(&proto.StreamResponse{
		Message: &proto.NotificationMessage{
			Info: "INFO2",
		},
	}); err != nil {
		log.Println("INFO2", err)
	}

	/*
		for {
			select {
			case e := <-ch:
				if err := stream.Send(&proto.StreamResponse{Message: e}); err != nil {
					log.Println("ERROR sending notification message over stream: ", err)
					return err
				}
			}
		}*/
	exit := make(chan bool)

	<-exit

	return nil
}
