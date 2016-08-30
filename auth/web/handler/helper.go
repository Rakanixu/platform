package handler

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"

	"golang.org/x/oauth2"
)

func SaveDatasource(url string, token *oauth2.Token) error {

	t := &go_micro_srv_datasource.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.String(),
	}
	c := go_micro_srv_datasource.NewDataSourceClient("go.micro.srv.desktop", nil)
	endpoint := &go_micro_srv_datasource.Endpoint{
		Url:   url,
		Token: t,
	}
	req := &go_micro_srv_datasource.CreateRequest{
		Endpoint: endpoint,
	}

	_, err := c.Create(context.TODO(), req)
	if err != nil {
		return err
	}
	return nil
}
