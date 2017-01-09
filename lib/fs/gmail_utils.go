package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (gfs *GmailFs) Authorize() (*datasource_proto.Token, error) {
	// Gmail Token is refreshed internally.
	return gfs.Endpoint.Token, nil
}
