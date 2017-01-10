package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (lfs *LocalFs) Authorize() (*datasource_proto.Token, error) {
	// No need to refresh
	return lfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (lfs *LocalFs) GetDatasourceId() string {
	return lfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (lfs *LocalFs) GetThumbnail(id string) (string, error) {
	return "", nil
}
