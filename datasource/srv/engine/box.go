package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"strings"
)

// Box struct
type Box struct {
	Endpoint proto.Endpoint
}

// Validate box datasource
func (b *Box) Validate(datasources string) (*proto.Endpoint, error) {
	str, err := globals.NewUUID()
	if err != nil {
		return &b.Endpoint, err
	}
	b.Endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	b.Endpoint.Id = globals.GetMD5Hash(b.Endpoint.Url + b.Endpoint.UserId)

	return &b.Endpoint, nil
}

// Save box data source
func (b *Box) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}
