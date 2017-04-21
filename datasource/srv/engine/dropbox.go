package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Dropbox struct
type Dropbox struct {
	Endpoint proto_datasource.Endpoint
}

// Validate dropbox datasource
func (s *Dropbox) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	s.Endpoint, err = GenerateEndpoint(ctx, s.Endpoint)
	if err != nil {
		return nil, err
	}

	return &s.Endpoint, nil
}

// Save dropbox data source
func (s *Dropbox) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete dropbox data source
func (s *Dropbox) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &s.Endpoint)
}

// CreateIndeWithAlias creates a index for dropbox datasource
func (s *Dropbox) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &s.Endpoint)
}
