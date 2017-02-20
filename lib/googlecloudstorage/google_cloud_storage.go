package googlecloudstorage

import (
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type GoogleCloudStorage struct{}

func NewGoogleCloudStorage() *GoogleCloudStorage {
	return &GoogleCloudStorage{}
}

var (
	SignedURLOption *SignedURLOptions
	GCSProjectID    string
)

// Loads the "default" google cloud config
func Register() {
	log.Println("INIT LIBRARY")

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

func (gcs *GoogleCloudStorage) CreateBucket(bucketName string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Insert(GCSProjectID, &storage.Bucket{
		Name: bucketName,
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

// Upload resource
func (gcs *GoogleCloudStorage) Upload(r io.Reader, bucketName, key string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Objects.Insert(bucketName, &storage.Object{
		Name: key,
	}).Media(r).Do()
	if err != nil {
		return err
	}

	return nil
}

// Download resource
func (gcs *GoogleCloudStorage) Download(bucketName, key string, opts ...string) (io.ReadCloser, error) {
	url, err := gcs.SignedObjectStorageURL(bucketName, key)
	if err != nil {
		return nil, err
	}

	cl := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Delete resource
func (gcs *GoogleCloudStorage) Delete(bucketName, key string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	if err := srv.Objects.Delete(bucketName, key).Do(); err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	// Refresh when resource expires. 1 hour availability since was requested.
	SignedURLOption.Expires = time.Now().Add(3600000000000)

	url, err := SignedURL(bucketName, objName, SignedURLOption)
	if err != nil {
		return "", err
	}

	return url, nil
}
