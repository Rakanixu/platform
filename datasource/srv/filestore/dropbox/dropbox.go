package dropbox

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Dropbox struct
type Dropbox struct {
	filestorer.FileStore
	Endpoint proto.Endpoint
}

// Validate slack
func (s *Dropbox) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &s.Endpoint, err
	}
	s.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	s.Endpoint.Id = globals.GetMD5Hash(s.Endpoint.Url + s.Endpoint.UserId)

	return &s.Endpoint, nil
}
