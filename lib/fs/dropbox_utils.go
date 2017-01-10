package fs

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	"net/url"
)

// Authorize
func (dfs *DropboxFs) Authorize() (*datasource_proto.Token, error) {
	// Dropbox Token never expires. No need to refresh.
	return dfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (dfs *DropboxFs) GetDatasourceId() string {
	return dfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (dfs *DropboxFs) GetThumbnail(id string) (string, error) {
	args := `{"path":"` + id + `","size":{".tag":"w640h480"}}`
	url := fmt.Sprintf("%s?authorization=%s&arg=%s", globals.DropboxThumbnailEndpoint, dfs.token(), url.QueryEscape(args))

	return url, nil
}
