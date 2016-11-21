package handler

import (
	"fmt"
	"log"

	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

//SaveDatasource call datasource-srv and save new data source
func SaveDatasource(ctx context.Context, user string, url string, token *oauth2.Token) error {

	c := proto_datasource.NewDataSourceClient("com.kazoup.srv.datasource", wrappers.NewKazoupClient())
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
		fmt.Print(err)
		return err
	}
	return nil
}

//PublishNotification send data source created notification
func PublishNotification(uID string) error {
	//c := client.NewClient()
	c := wrappers.NewKazoupClient()
	n := &notification_proto.NotificationMessage{
		Info:   "Datasource created succesfully",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: string(uID),
	}

	// Publish scan topic, crawlers should pick up message and start scanning
	return c.Publish(globals.NewSystemContext(), c.NewPublication(globals.NotificationTopic, n))
}
