package handler

import (
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

type Service struct{}

func (s *Service) Stream(ctx context.Context, stream server.Streamer) error {
	// Listen for StreamRequest (this is blocking)
	req := &proto_notification.StreamRequest{}
	if err := stream.Recv(req); err != nil {
		return err
	}

	// StreamNotifications subscribes to NotificationTopic and return channels for communications
	ch, exit, err := StreamNotifications(ctx, req)
	if err != nil {
		return err
	}

	defer func() {
		close(exit)
		stream.Close()
	}()

	for {
		select {
		// Listen over the open channel, all received notification will be pushed over this channel
		// Once channel retrieves data, send it back over the stream
		case e := <-ch:
			if err := stream.Send(&proto_notification.StreamResponse{Message: e}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_notification.HealthRequest, rsp *proto_notification.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
