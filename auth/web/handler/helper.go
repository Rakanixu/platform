package handler

import (
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func SaveDatasource(ctx context.Context, user string, url string, token *oauth2.Token) error {
	c := proto_datasource.NewDataSourceClient(globals.DATASOURCE_SERVICE_NAME, nil)

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
