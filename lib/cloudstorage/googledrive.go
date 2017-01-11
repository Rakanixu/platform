package cloudstorage

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"io"
	"time"
)

type GoogleDriveCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewGoogleDriveCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &GoogleDriveCloudStorage{
		Endpoint: e,
	}
}

// Upload
func (gcs *GoogleDriveCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download
func (gcs *GoogleDriveCloudStorage) Download(fileID string, opts ...string) (io.ReadCloser, error) {
	srv, err := gcs.getDriveService()
	if err != nil {
		return nil, err
	}

	res, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

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
