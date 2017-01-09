package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (gfs *GoogleDriveFs) Authorize() (*datasource_proto.Token, error) {
	// GoogleDrive Token is refreshed internally.
	return gfs.Endpoint.Token, nil
}
