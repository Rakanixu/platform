package fs

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/oauth2"
	"time"
)

// Authorize
func (bfs *BoxFs) Authorize() (*datasource_proto.Token, error) {
	tokenSource := globals.NewBoxOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  bfs.Endpoint.Token.AccessToken,
		TokenType:    bfs.Endpoint.Token.TokenType,
		RefreshToken: bfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(bfs.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	bfs.Endpoint.Token.AccessToken = t.AccessToken
	bfs.Endpoint.Token.TokenType = t.TokenType
	bfs.Endpoint.Token.RefreshToken = t.RefreshToken
	bfs.Endpoint.Token.Expiry = t.Expiry.Unix()

	return bfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (bfs *BoxFs) GetDatasourceId() string {
	return bfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to an image
func (bfs *BoxFs) GetThumbnail(id string) (string, error) {
	url := fmt.Sprintf(
		"%s%s&Authorization=%s",
		globals.BoxFileMetadataEndpoint,
		id,
		"/thumbnail.png?min_height=256&min_width=256",
		bfs.token(),
	)

	return url, nil
}
