package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/handler"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	l := int64(len(handler.Crawls)) + 1

	s, err := handler.MapScanner(l, endpoint)
	if err != nil {
		return err
	}

	handler.Crawls[l] = s
	s.Start(handler.Crawls, l)

	// Publish notification
	msg := &notification_proto.NotificationMessage{
		Info: "Scan started on " + endpoint.Url + " datasource.",
	}

	if err := client.Publish(ctx, client.NewPublication(globals.NotificationTopic, msg)); err != nil {
		return err
	}

	return nil
}
