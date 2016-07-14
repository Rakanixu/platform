package local

import (
	"os"
	"strings"

	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
)

// Local struct
type Local struct {
	filestorer.FileStore
	Endpoint   proto.Endpoint
	DataOrigin string
}

// Validate local datasource (directory exists)
func (l *Local) Validate() error {
	i := strings.LastIndex(l.Endpoint.Url, "//")

	l.DataOrigin = l.Endpoint.Url[i+1 : len(l.Endpoint.Url)] // Local filesystem path
	if _, err := os.Stat(l.DataOrigin); os.IsNotExist(err) {
		return err
	}

	return nil
}
