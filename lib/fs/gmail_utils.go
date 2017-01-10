package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (gfs *GmailFs) Authorize() (*datasource_proto.Token, error) {
	// Gmail Token is refreshed internally.
	return gfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (gfs *GmailFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (gfs *GmailFs) GetThumbnail(id string) (string, error) {
	return "", nil
}
