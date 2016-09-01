package onedrive

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Fake struct
type Onedrive struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Onedrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	return nil, nil
}
