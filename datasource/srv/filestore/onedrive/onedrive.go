package onedrive

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
)

// Fake struct
type Onedrive struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Onedrive) Validate() error {
	return nil
}
