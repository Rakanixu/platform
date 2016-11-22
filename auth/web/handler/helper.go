package handler

import (
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/wrappers"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

//SaveDatasource call datasource-srv and save new data source
func SaveDatasource(ctx context.Context, user string, url string, token *oauth2.Token) error {
	client := wrappers.NewKazoupClient()

	srvReq := client.NewRequest(
		globals.DATASOURCE_SERVICE_NAME,
		"DataSource.Create",
		&proto_datasource.CreateRequest{
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
		},
	)

	srvRes := &proto_datasource.CreateResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.New("auth.web.helper.DataSource.Create", err.Error(), 500)
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
