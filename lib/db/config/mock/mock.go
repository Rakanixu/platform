package mock

import (
	"github.com/kazoup/platform/lib/db/config"
	"github.com/kazoup/platform/lib/db/config/proto/config"
	"golang.org/x/net/context"
)

type mock struct{}

func init() {
	config.Register(new(mock))
}

func (e *mock) Init() error {
	return nil
}

func (e *mock) CreateIndex(ctx context.Context, req *proto_config.CreateIndexRequest) (*proto_config.CreateIndexResponse, error) {
	return &proto_config.CreateIndexResponse{}, nil
}

func (e *mock) Status(ctx context.Context, req *proto_config.StatusRequest) (*proto_config.StatusResponse, error) {
	return &proto_config.StatusResponse{}, nil
}

func (e *mock) AddAlias(ctx context.Context, req *proto_config.AddAliasRequest) (*proto_config.AddAliasResponse, error) {
	return &proto_config.AddAliasResponse{}, nil
}

func (e *mock) DeleteIndex(ctx context.Context, req *proto_config.DeleteIndexRequest) (*proto_config.DeleteIndexResponse, error) {
	return &proto_config.DeleteIndexResponse{}, nil
}

func (e *mock) DeleteAlias(ctx context.Context, req *proto_config.DeleteAliasRequest) (*proto_config.DeleteAliasResponse, error) {
	return &proto_config.DeleteAliasResponse{}, nil
}

func (e *mock) RenameAlias(ctx context.Context, req *proto_config.RenameAliasRequest) (*proto_config.RenameAliasResponse, error) {
	return &proto_config.RenameAliasResponse{}, nil
}
