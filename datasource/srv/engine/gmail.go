package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Gmail struct
type Gmail struct {
	Endpoint datasource_proto.Endpoint
}

// Validate gmail data  source
func (g *Gmail) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	s, err := globals.NewUUID()
	if err != nil {
		return &g.Endpoint, err
	}
	g.Endpoint.Index = "index" + strings.Replace(s, "-", "", 1)
	g.Endpoint.Id = globals.GetMD5Hash(g.Endpoint.Url + g.Endpoint.UserId)

	return &g.Endpoint, nil
}

// Save gmail data source
func (g *Gmail) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete gmail data source
func (g *Gmail) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &g.Endpoint)
}
