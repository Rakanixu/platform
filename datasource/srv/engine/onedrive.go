package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Onedrive struct
type Onedrive struct {
	Endpoint proto_datasource.Endpoint
}

// Validate
func (o *Onedrive) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	o.Endpoint, err = GenerateEndpoint(ctx, o.Endpoint)
	if err != nil {
		return nil, err
	}

	return &o.Endpoint, nil
}

// Save one drive datasource
func (o *Onedrive) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete one drive data source
func (o *Onedrive) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &o.Endpoint)
}

// CreateIndeWithAlias creates a index for local datasource
func (o *Onedrive) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &o.Endpoint)
}
