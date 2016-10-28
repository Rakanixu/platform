package engine

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Gmail struct
type Gmail struct {
	Endpoint datasource_proto.Endpoint
}

// Validate gmail data  source
func (g *Gmail) Validate(datasources string) (*datasource_proto.Endpoint, error) {
	return GenerateEndpoint(&g.Endpoint)
}

// Save gmail data source
func (g *Gmail) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete gmail data source
func (g *Gmail) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &g.Endpoint)
}

// Scan gmail data source
func (g *Gmail) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &g.Endpoint)
}

// CreateIndeWithAlias creates a index for gmail datasource
func (g *Gmail) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &g.Endpoint)
}
