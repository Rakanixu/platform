package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"golang.org/x/net/context"
)

type Mock struct {
	Endpoint proto_datasource.Endpoint
}

func (m *Mock) Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error) {
	return &m.Endpoint, nil
}

func (m *Mock) Save(ctx context.Context, data interface{}, id string) error {
	return nil
}

func (m *Mock) Delete(ctx context.Context) error {
	return nil
}

func (m *Mock) CreateIndexWithAlias(ctx context.Context) error {
	return nil
}
