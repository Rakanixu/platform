package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (dfs *DropboxFs) Authorize() (*datasource_proto.Token, error) {
	// Dropbox Token never expires. No need to refresh.
	return dfs.Endpoint.Token, nil
}
