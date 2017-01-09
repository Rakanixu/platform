package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (lfs *LocalFs) Authorize() (*datasource_proto.Token, error) {
	// No need to refresh
	return lfs.Endpoint.Token, nil
}
