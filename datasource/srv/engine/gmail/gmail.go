package gmail

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Gmail struct
type Gmail struct {
	engine.DataSource
	Endpoint datasource_proto.Endpoint
}

// Validate
func (g *Gmail) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	s, err := globals.NewUUID()
	if err != nil {
		return &g.Endpoint, err
	}
	g.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	g.Endpoint.Id = globals.GetMD5Hash(g.Endpoint.Url + g.Endpoint.UserId)

	return &g.Endpoint, nil
}
