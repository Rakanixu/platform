package fs

import (
	"encoding/json"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/onedrive"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

// Authorize
func (ofs *OneDriveFs) Authorize() (*datasource_proto.Token, error) {
	tokenSource := globals.NewMicrosoftOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  ofs.Endpoint.Token.AccessToken,
		TokenType:    ofs.Endpoint.Token.TokenType,
		RefreshToken: ofs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(ofs.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	ofs.Endpoint.Token.AccessToken = t.AccessToken
	ofs.Endpoint.Token.TokenType = t.TokenType
	ofs.Endpoint.Token.RefreshToken = t.RefreshToken
	ofs.Endpoint.Token.Expiry = t.Expiry.Unix()

	return ofs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (ofs *OneDriveFs) GetDatasourceId() string {
	return ofs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (ofs *OneDriveFs) GetThumbnail(id string) (string, error) {
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := fmt.Sprintf("%sitems/%s/thumbnails/0/medium", globals.OneDriveEndpoint+Drive, id)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.token())
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var thumbRsp *onedrive.FileThumbnailResponse
	if err := json.NewDecoder(res.Body).Decode(&thumbRsp); err != nil {
		return "", err
	}

	return thumbRsp.URL, nil
}
