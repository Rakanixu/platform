package cloudstorage

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
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
}

// CloudStorageUtils interface
type CloudStorageUtils interface {
	SignedObjectStorageURL(string, string) (string, error)
	DeleteBucket() error
}

// NewCloudStorageFromEndpoint constructor
func NewCloudStorageFromEndpoint(e *datasource_proto.Endpoint, connector string) (CloudStorage, error) {
	switch connector {
	case globals.GoogleCloudStorage:
		return NewGoogleCloudStorage(e), nil
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

	return nil, errors.New("Error parsing URL")
}
