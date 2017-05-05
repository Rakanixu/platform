package objectstorage

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/objectstorage"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type GoogleCloudStorage struct{}

var (
	SignedURLOption *SignedURLOptions
	GCSProjectID    string
)

func init() {
	objectstorage.Register(new(GoogleCloudStorage))
}

// Loads the "default" google cloud config
func (gcs *GoogleCloudStorage) Init() error {
	b, err := ioutil.ReadFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	if err != nil {
		return err
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
		return err
	}

	GCSProjectID = d.ProjectID

	SignedURLOption = &SignedURLOptions{
		GoogleAccessID: d.ClientEmail,
		PrivateKey:     []byte(d.PrivateKey),
		Method:         http.MethodGet,
	}

	return nil
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

// DeleteBucket deletes a bucket and all its contents in GCS account
func (gcs *GoogleCloudStorage) DeleteBucket(bucketName string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Get(bucketName).Do()
	// Bucket does not exists
	if err != nil {
		return nil // If does not exists, it's already deleted, then success
	}

	r, err := srv.Objects.List(bucketName).Fields("items,nextPageToken").Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(bucketName, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gcs.deleteBucketNextPage(srv, bucketName, r.NextPageToken); err != nil {
			return err
		}
	}

	if err := srv.Buckets.Delete(bucketName).Do(); err != nil {
		return err
	}

	return nil
}

// Upload resource
func (gcs *GoogleCloudStorage) Upload(r io.ReadCloser, bucketName, key string) error {
	defer r.Close()

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

func (gcs *GoogleCloudStorage) deleteBucketNextPage(srv *storage.Service, bucketName, nextPageToken string) error {
	r, err := srv.Objects.List(bucketName).Fields("items,nextPageToken").PageToken(nextPageToken).Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(bucketName, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gcs.deleteBucketNextPage(srv, bucketName, r.NextPageToken); err != nil {
			return err
		}
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
