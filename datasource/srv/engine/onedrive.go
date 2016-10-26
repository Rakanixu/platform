package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Onedrive struct
type Onedrive struct {
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

// Save one drive datasource
func (o *Onedrive) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete one drive data source
func (o *Onedrive) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &o.Endpoint)
}
