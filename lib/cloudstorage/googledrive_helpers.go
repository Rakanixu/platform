package cloudstorage

import (
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"time"
)

// getDriveService return a google drive service instance
func (gcs *GoogleDriveCloudStorage) getDriveService() (*drive.Service, error) {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gcs.Endpoint.Token.AccessToken,
		TokenType:    gcs.Endpoint.Token.TokenType,
		RefreshToken: gcs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gcs.Endpoint.Token.Expiry, 0),
	})

	return drive.New(c)
}
