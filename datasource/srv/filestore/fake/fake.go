package fake

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Fake struct
type Fake struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (l *Fake) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	return nil, nil
}
