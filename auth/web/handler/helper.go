package handler

import (
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func SaveDatasource(url string, token *oauth2.Token) error {
	c := proto_datasource.NewDataSourceClient("", nil)

	req := &proto_datasource.CreateRequest{
		Endpoint: &proto_datasource.Endpoint{
			Url: url,
			Token: &proto_datasource.Token{
				AccessToken:  token.AccessToken,
				TokenType:    token.TokenType,
				RefreshToken: token.RefreshToken,
				Expiry:       token.Expiry.String(),
			},
		},
	}

	_, err := c.Create(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
