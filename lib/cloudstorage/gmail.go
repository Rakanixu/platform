package cloudstorage

import (
	"bytes"
	"encoding/base64"
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
	"io"
	"io/ioutil"
	"time"
)

type GmailCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewGmailCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &GmailCloudStorage{
		Endpoint: e,
	}
}

// Upload
func (gcs *GmailCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download
func (gcs *GmailCloudStorage) Download(fileID string, opts ...string) (io.ReadCloser, error) {
	cfg := globals.NewGmailOauthConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gcs.Endpoint.Token.AccessToken,
		TokenType:    gcs.Endpoint.Token.TokenType,
		RefreshToken: gcs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gcs.Endpoint.Token.Expiry, 0),
	})

	s, err := gmail.New(c)
	if err != nil {
		return nil, err
	}

	srv := gmail.NewUsersMessagesService(s)

	if len(opts) == 0 {
		return nil, errors.New("ERROR opts required (attachmentID)")
	}
	mpb, err := srv.Attachments.Get("me", fileID, opts[0]).Fields("data").Do()
	if err != nil {
		return nil, err
	}

	b, err := base64.URLEncoding.DecodeString(mpb.Data)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil
}
