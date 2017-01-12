package cloudstorage

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io"
	"net/http"
)

type DropboxCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewDropboxCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &DropboxCloudStorage{
		Endpoint: e,
	}
}

// Upload resource
func (dcs *DropboxCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download resource
func (dcs *DropboxCloudStorage) Download(fileID string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, globals.DropboxFileDownload, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dcs.token())
	req.Header.Set("Dropbox-API-Arg", `{
			"path": "`+fileID+`"
		}`)
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

// Delete resource
func (dcs *DropboxCloudStorage) Delete(bucketName string, objName string) error {
	return nil
}
