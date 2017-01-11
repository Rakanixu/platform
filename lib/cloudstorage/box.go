package cloudstorage

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"io"
	"net/http"
)

type BoxCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewBoxCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &BoxCloudStorage{
		Endpoint: e,
	}
}

// Upload resource
func (bcs *BoxCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download resource
func (bcs *BoxCloudStorage) Download(fileID string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/content", globals.BoxFileMetadataEndpoint, fileID)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", bcs.token())
	rsp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

// Delete resource
func (bcs *BoxCloudStorage) Delete(bucketName string, objName string) error {
	return nil
}

func (bcs *BoxCloudStorage) token() string {
	return "Bearer " + bcs.Endpoint.Token.AccessToken
}
