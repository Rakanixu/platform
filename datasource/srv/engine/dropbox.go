package engine

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Dropbox struct
type Dropbox struct {
	Endpoint proto.Endpoint
}

// Validate dropbox datasource
func (s *Dropbox) Validate(ctx context.Context, c client.Client, datasources string) (*proto.Endpoint, error) {
	var err error

	s.Endpoint, err = GenerateEndpoint(ctx, c, s.Endpoint)
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
func (s *Dropbox) Delete(ctx context.Context, c client.Client) error {
	return DeleteDataSource(ctx, c, &s.Endpoint)
}

// Scan dropbox data source
func (s *Dropbox) Scan(ctx context.Context, c client.Client) error {
	return ScanDataSource(ctx, c, &s.Endpoint)
}

// CreateIndeWithAlias creates a index for dropbox datasource
func (s *Dropbox) CreateIndexWithAlias(ctx context.Context, c client.Client) error {
	return CreateIndexWithAlias(ctx, c, &s.Endpoint)
}
