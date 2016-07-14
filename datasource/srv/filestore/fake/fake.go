package fake

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
)

// Fake struct
type Fake struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (l *Fake) Validate() error {
	return nil
}
