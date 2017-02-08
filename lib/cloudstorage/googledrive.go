package cloudstorage

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"io"
	"net/http"
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

	var res *http.Response
	if len(opts) > 0 {
		// Google documents are special, so the have to be exported
		if opts[0] == "export" {
			fmt.Println(opts[1])

			res, err = srv.Files.Export(fileID, opts[1]).Download()
			if err != nil {
				return nil, err
			}
		}

		if opts[0] == "download" {
			res, err = srv.Files.Get(fileID).Download()
			if err != nil {
				return nil, err
			}
		}
	}

	return res.Body, nil
}

// Delete resource
func (gcs *GoogleDriveCloudStorage) Delete(bucketName string, objName string) error {
	return nil
}
