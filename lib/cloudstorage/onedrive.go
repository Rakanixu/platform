package cloudstorage

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io"
	"net/http"
)

type OneDriveCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewOneDriveCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &OneDriveCloudStorage{
		Endpoint: e,
	}
}

// Upload
func (ocs *OneDriveCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download
func (ocs *OneDriveCloudStorage) Download(fileID string, opts ...string) (io.ReadCloser, error) {
	oc := &http.Client{}
	// https://dev.onedrive.com/items/download.htm
	url := globals.OneDriveEndpoint + "drive/items/" + fileID + "/content"
	oreq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	oreq.Header.Set("Authorization", ocs.token())
	res, err := oc.Do(oreq)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Delete resource
func (ocs *OneDriveCloudStorage) Delete(bucketName string, objName string) error {
	return nil
}
