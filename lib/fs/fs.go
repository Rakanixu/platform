package fs

import (
	"encoding/json"
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Fs File System interface. Fyle system is responsible to manage its own files
type Fs interface {
	List(client.Client) (chan file.File, chan bool, error)
	CreateFile(context.Context, client.Client, file_proto.CreateRequest) (*file_proto.CreateResponse, error)
	DeleteFile(context.Context, client.Client, file_proto.DeleteRequest) (*file_proto.DeleteResponse, error)
	ShareFile(context.Context, client.Client, file_proto.ShareRequest) (string, error)
	DownloadFile(string, client.Client, ...string) (io.ReadCloser, error)
	UploadFile(io.Reader, string) error
	SignedObjectStorageURL(string) (string, error)
	DeleteIndexBucketFromGCS() error
	GetDatasourceId() string
	Token(client.Client) string
	GetThumbnail(string, client.Client) (string, error)
}

var (
	SignedURLOption *SignedURLOptions
	GCSProjectID    string
)

// Loads the google cloud config
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

// NewFsFromEndpoint constructor from endpoint
func NewFsFromEndpoint(e *datasource_proto.Endpoint) (Fs, error) {
	dsUrl := strings.Split(e.Url, ":")

	switch dsUrl[0] {
	case globals.Local:
		return NewLocalFsFromEndpoint(e), nil
	case globals.Slack:
		return NewSlackFsFromEndpoint(e), nil
	case globals.GoogleDrive:
		return NewGoogleDriveFsFromEndpoint(e), nil
	case globals.Gmail:
		return NewGmailFsFromEndpoint(e), nil
	case globals.OneDrive:
		return NewOneDriveFsFromEndpoint(e), nil
	case globals.Dropbox:
		return NewDropboxFsFromEndpoint(e), nil
	case globals.Box:
		return NewBoxFsFromEndpoint(e), nil
	default:
		return nil, errors.New("Not such file system (fs)")
	}

	return nil, errors.New("Error parsing URL")
}

//SignedObjectStorageURL returns a temporary valid URL to the resource
func SignedObjectStorageURL(index, objName string) (string, error) {
	// Refresh when resource expires. 1 hour availability since was requested.
	SignedURLOption.Expires = time.Now().Add(3600000000000)

	url, err := SignedURL(index, objName, SignedURLOption)
	if err != nil {
		return "", err
	}

	return url, nil
}

// UploadFile uploads a file into google cloud storage
func UploadFile(r io.Reader, index, fID string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Get(index).Do()
	// Bucket does not exists, create it
	if err != nil {
		_, err = srv.Buckets.Insert(GCSProjectID, &storage.Bucket{
			Name: index,
		}).Do()
		if err != nil {
			return err
		}
	}

	_, err = srv.Objects.Insert(index, &storage.Object{
		Name: fID,
	}).Media(r).Do()
	if err != nil {
		return err
	}

	return nil
}

// DeleteFile deletes a file from google cloud storage
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

// DeleteBucket deletes a bucket and all its contents
func DeleteBucket(index, nextPageToken string) error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Get(index).Do()
	// Bucket does not exists
	if err != nil {
		return nil // If does not exists, it's already deleted, then success
	}

	r, err := srv.Objects.List(index).Fields("items,nextPageToken").Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(index, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := deleteBucketNextPage(srv, index, r.NextPageToken); err != nil {
			return err
		}
	}

	if err := srv.Buckets.Delete(index).Do(); err != nil {
		return err
	}

	return nil
}

func deleteBucketNextPage(srv *storage.Service, index, nextPageToken string) error {
	r, err := srv.Objects.List(index).Fields("items,nextPageToken").PageToken(nextPageToken).Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(index, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := deleteBucketNextPage(srv, index, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}
