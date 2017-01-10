package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (sfs *SlackFs) Authorize() (*datasource_proto.Token, error) {
	// Gmail Token is refreshed internally.
	return sfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (sfs *SlackFs) GetDatasourceId() string {
	return sfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (sfs *SlackFs) GetThumbnail(id string) (string, error) {
	return "", nil
}
