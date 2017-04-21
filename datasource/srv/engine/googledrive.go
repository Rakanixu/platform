package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Googledrive struct
type Googledrive struct {
	Endpoint proto_datasource.Endpoint
}

// Validate google drive data source
func (g *Googledrive) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	g.Endpoint, err = GenerateEndpoint(ctx, g.Endpoint)
	if err != nil {
		return nil, err
	}

	return &g.Endpoint, nil
}

// Save google drive data source
func (g *Googledrive) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete google drive data source
func (g *Googledrive) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &g.Endpoint)
}

// CreateIndeWithAlias creates a index for google drive datasource
func (g *Googledrive) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &g.Endpoint)
}
