package cloudstorage

import (
	"errors"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io"
)

// CloudStorage interface
type CloudStorage interface {
	CloudStorageOperations
	CloudStorageUtils
}

// CloudStorageOperations interface
type CloudStorageOperations interface {
	Upload(io.Reader, string) error
	Download(string, ...string) (io.ReadCloser, error)
	Delete(string, string) error
}

// CloudStorageUtils interface
type CloudStorageUtils interface {
	SignedObjectStorageURL(string, string) (string, error)
	DeleteBucket() error
}

// NewCloudStorageFromEndpoint constructor
func NewCloudStorageFromEndpoint(e *proto_datasource.Endpoint, connector string) (CloudStorage, error) {
	switch connector {
	case globals.Slack:
		return NewSlackCloudStorage(e), nil
	case globals.GoogleDrive:
		return NewGoogleDriveCloudStorage(e), nil
	case globals.Gmail:
		return NewGmailCloudStorage(e), nil
	case globals.OneDrive:
		return NewOneDriveCloudStorage(e), nil
	case globals.Dropbox:
		return NewDropboxCloudStorage(e), nil
	case globals.Box:
		return NewBoxCloudStorage(e), nil
	default:
		return nil, errors.New("Not such cloud storage constructor")
	}

}
