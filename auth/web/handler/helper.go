package handler

import (
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
)

func SaveDatasource(ctx context.Context, user string, url string, token *oauth2.Token) error {

	c := proto_datasource.NewDataSourceClient(globals.DATASOURCE_SERVICE_NAME, wrappers.NewKazoupClient())
	log.Print(c)
	req := &proto_datasource.CreateRequest{
		Endpoint: &proto_datasource.Endpoint{
			UserId:          user,
			Url:             url,
			LastScan:        0,
			LastScanStarted: 0,
			CrawlerRunning:  false,
			Token: &proto_datasource.Token{
				AccessToken:  token.AccessToken,
				TokenType:    token.TokenType,
				RefreshToken: token.RefreshToken,
				Expiry:       token.Expiry.Unix(),
			},
		},
	}

	_, err := c.Create(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func PublishNotification(uID string) error {
	c := client.NewClient()
	n := &notification_proto.NotificationMessage{
		Info:   "Datasource created succesfully",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: string(uID),
	}

	// Publish scan topic, crawlers should pick up message and start scanning
	return c.Publish(globals.NewSystemContext(), c.NewPublication(globals.NotificationTopic, n))
}
