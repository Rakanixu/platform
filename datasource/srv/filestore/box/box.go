package box

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Dropbox struct
type Box struct {
	filestorer.FileStore
	Endpoint proto.Endpoint
}

// Validate slack
func (b *Box) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &b.Endpoint, err
	}
	b.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	b.Endpoint.Id = globals.GetMD5Hash(b.Endpoint.Url + b.Endpoint.UserId)

	return &b.Endpoint, nil
}
