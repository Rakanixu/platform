package fs

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Authorize
func (gfs *GoogleDriveFs) Authorize() (*datasource_proto.Token, error) {
	// GoogleDrive Token is refreshed internally.
	return gfs.Endpoint.Token, nil
}

// GetDatasourceId returns datasource ID
func (gfs *GoogleDriveFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (gfs *GoogleDriveFs) GetThumbnail(id string) (string, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return "", err
	}
	r, err := srv.Files.Get(id).Fields("thumbnailLink").Do()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%ss700", r.ThumbnailLink[:len(r.ThumbnailLink)-4]), nil
}
