package gmail

import (
	filestorer "github.com/kazoup/platform/datasource/srv/filestore"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Fake struct
type Gmail struct {
	Endpoint datasource_proto.Endpoint
	filestorer.FileStore
}

// Validate fake, always fine
func (g *Gmail) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	s, err := globals.NewUUID()
	if err != nil {
		return &g.Endpoint, err
	}
	g.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	g.Endpoint.Id = globals.GetMD5Hash(g.Endpoint.Url + g.Endpoint.UserId)

	return &g.Endpoint, nil
}
