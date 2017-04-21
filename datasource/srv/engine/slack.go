package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

// Slack struct
type Slack struct {
	Endpoint proto_datasource.Endpoint
}

// Validate slack data source
func (s *Slack) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	var err error

	s.Endpoint, err = GenerateEndpoint(ctx, s.Endpoint)
	if err != nil {
		return nil, err
	}

	return &s.Endpoint, nil
}

// Save slack datasource
func (s *Slack) Save(ctx context.Context, data interface{}, id string) error {
	return SaveDataSource(ctx, data, id)
}

// Delete slack data source
func (s *Slack) Delete(ctx context.Context) error {
	return DeleteDataSource(ctx, &s.Endpoint)
}

// CreateIndeWithAlias creates a index for local datasource
func (s *Slack) CreateIndexWithAlias(ctx context.Context) error {
	return CreateIndexWithAlias(ctx, &s.Endpoint)
}
