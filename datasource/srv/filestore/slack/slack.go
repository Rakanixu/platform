package slack

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
)

// Fake struct
type Slack struct {
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Slack) Validate() error {
	return nil
}
