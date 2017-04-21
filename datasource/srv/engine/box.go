package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Box struct
type Box struct {
	Endpoint proto_datasource.Endpoint
}

// Validate box datasource
func (b *Box) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	b.Endpoint, err = GenerateEndpoint(ctx, b.Endpoint)
	if err != nil {
		return nil, err
	}

	return &b.Endpoint, nil
}

// Save box data source
func (b *Box) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete box data source
func (b *Box) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &b.Endpoint)
}

// CreateIndeWithAlias creates a index for box datasource
func (b *Box) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &b.Endpoint)
}
