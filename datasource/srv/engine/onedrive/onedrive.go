package onedrive

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Onedrive struct
type Onedrive struct {
	engine.DataSource
	Endpoint datasource_proto.Endpoint
}

// Validate
func (o *Onedrive) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	if len(o.Endpoint.Index) == 0 {
		s, err := globals.NewUUID()
		if err != nil {
			return &o.Endpoint, err
		}
		o.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	}
	o.Endpoint.Id = globals.GetMD5Hash(o.Endpoint.Url + o.Endpoint.UserId)

	return &o.Endpoint, nil
}
