package cloudstorage

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"io"
	"net/http"
)

type SlackCloudStorage struct {
	Endpoint *datasource_proto.Endpoint
}

// NewBoxCloudStorage
func NewSlackCloudStorage(e *datasource_proto.Endpoint) CloudStorage {
	return &SlackCloudStorage{
		Endpoint: e,
	}
}

// Upload
func (scs *SlackCloudStorage) Upload(r io.Reader, fileID string) error {
	return nil
}

// Download
func (scs *SlackCloudStorage) Download(url string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", scs.token())
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// Delete resource
func (scs *SlackCloudStorage) Delete(bucketName string, objName string) error {
	return nil
}
