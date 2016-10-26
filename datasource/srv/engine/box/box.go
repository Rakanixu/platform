package box

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

// Box struct
type Box struct {
	engine.DataSource
	Endpoint proto.Endpoint
}

// Validate
func (b *Box) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &b.Endpoint, err
	}
	b.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	b.Endpoint.Id = globals.GetMD5Hash(b.Endpoint.Url + b.Endpoint.UserId)

	return &b.Endpoint, nil
}
