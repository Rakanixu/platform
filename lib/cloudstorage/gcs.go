package cloudstorage

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GoogleCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

func NewGoogleCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &GoogleCloudStorage{
		Endpoint: e,
	}
}

var (
	SignedURLOption *SignedURLOptions
	GCSProjectID    string
)

// Loads the "default" google cloud config
func init() {
	b, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var d struct {
		// Common fields
		Type      string
		ProjectID string `json:"project_id"`
		ClientID  string `json:"client_id"`

		// User Credential fields
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`

		// Service Account fields
		ClientEmail  string `json:"client_email"`
		PrivateKeyID string `json:"private_key_id"`
		PrivateKey   string `json:"private_key"`
	}

	if err := json.Unmarshal(b, &d); err != nil {
		log.Fatalf("Error unmarshalling config file: %s", err)
	}

	GCSProjectID = d.ProjectID

	SignedURLOption = &SignedURLOptions{
		GoogleAccessID: d.ClientEmail,
		PrivateKey:     []byte(d.PrivateKey),
		Method:         http.MethodGet,
	}
}

func (gcs *GoogleCloudStorage) Upload(r io.Reader, fileID string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Get(gcs.Endpoint.Index).Do()
	// Bucket does not exists, create it
	if err != nil {
		_, err = srv.Buckets.Insert(GCSProjectID, &storage.Bucket{
			Name: gcs.Endpoint.Index,
		}).Do()
		if err != nil {
			return err
		}
	}

	_, err = srv.Objects.Insert(gcs.Endpoint.Index, &storage.Object{
		Name: fileID,
	}).Media(r).Do()
	if err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) Download(string, ...string) (io.ReadCloser, error) {
	return nil, nil
}

// DeleteFile deletes a file from google cloud storage
/*
func DeleteFile(index, fID string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	if err := srv.Objects.Delete(index, fID).Do(); err != nil {
		return err
	}

	return nil
}
*/
