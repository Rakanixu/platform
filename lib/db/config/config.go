package config

import (
	"github.com/kazoup/platform/lib/db/config/proto/config"
	"golang.org/x/net/context"
)

type DBConfig interface {
	Init() error
	Config
}

type Config interface {
	CreateIndex(ctx context.Context, req *proto_config.CreateIndexRequest) (*proto_config.CreateIndexResponse, error)
	Status(ctx context.Context, req *proto_config.StatusRequest) (*proto_config.StatusResponse, error)
	AddAlias(ctx context.Context, req *proto_config.AddAliasRequest) (*proto_config.AddAliasResponse, error)
	DeleteIndex(ctx context.Context, req *proto_config.DeleteIndexRequest) (*proto_config.DeleteIndexResponse, error)
	DeleteAlias(ctx context.Context, req *proto_config.DeleteAliasRequest) (*proto_config.DeleteAliasResponse, error)
	RenameAlias(ctx context.Context, req *proto_config.RenameAliasRequest) (*proto_config.RenameAliasResponse, error)
}

var (
	config DBConfig
)

func Register(storage DBConfig) {
	config = storage
}

func Init() error {
	return config.Init()
}

func CreateIndex(ctx context.Context, req *proto_config.CreateIndexRequest) (*proto_config.CreateIndexResponse, error) {
	return config.CreateIndex(ctx, req)
}

func Status(ctx context.Context, req *proto_config.StatusRequest) (*proto_config.StatusResponse, error) {
	return config.Status(ctx, req)
}

func AddAlias(ctx context.Context, req *proto_config.AddAliasRequest) (*proto_config.AddAliasResponse, error) {
	return config.AddAlias(ctx, req)
}

func DeleteIndex(ctx context.Context, req *proto_config.DeleteIndexRequest) (*proto_config.DeleteIndexResponse, error) {
	return config.DeleteIndex(ctx, req)
}

func DeleteAlias(ctx context.Context, req *proto_config.DeleteAliasRequest) (*proto_config.DeleteAliasResponse, error) {
	return config.DeleteAlias(ctx, req)
}

func RenameAlias(ctx context.Context, req *proto_config.RenameAliasRequest) (*proto_config.RenameAliasResponse, error) {
	return config.RenameAlias(ctx, req)
}
