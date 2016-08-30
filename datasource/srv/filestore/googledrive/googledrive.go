package googledrive

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
)

// Fake struct
type Googledrive struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Googledrive) Validate() error {
	return nil
}
