package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Gmail struct
type Gmail struct {
	Endpoint proto_datasource.Endpoint
}

// Validate gmail data  source
func (g *Gmail) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	g.Endpoint, err = GenerateEndpoint(ctx, g.Endpoint)
	if err != nil {
		return nil, err
	}

	return &g.Endpoint, nil
}

// Save gmail data source
func (g *Gmail) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete gmail data source
func (g *Gmail) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &g.Endpoint)
}

// CreateIndeWithAlias creates a index for gmail datasource
func (g *Gmail) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &g.Endpoint)
}
