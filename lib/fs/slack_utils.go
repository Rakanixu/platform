package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (sfs *SlackFs) Authorize() (*datasource_proto.Token, error) {
	// Gmail Token is refreshed internally.
	return sfs.Endpoint.Token, nil
}
