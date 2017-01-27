package handler

import (
	engine "github.com/kazoup/platform/db/srv/engine"
	proto "github.com/kazoup/platform/db/srv/proto/config"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Config struct
type Config struct{}

// CreateIndexWithSettings Config srv handler
func (cf *Config) CreateIndex(ctx context.Context, req *proto.CreateIndexRequest, rsp *proto.CreateIndexResponse) error {
	_, err := engine.CreateIndex(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.CreateIndexWithSettings", err.Error())
	}

	return nil
}

// Status Config srv handler
func (cf *Config) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	response, err := engine.Status(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Status", err.Error())
	}

	rsp.Status = response.Status

	return nil
}

// RenameIndexAlias Config srv handler
func (cf *Config) AddAlias(ctx context.Context, req *proto.AddAliasRequest, rsp *proto.AddAliasResponse) error {
	_, err := engine.AddAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.AddAlias", err.Error())
	}

	return nil
}

// DeleteIndex Config srv handler
func (cf *Config) DeleteIndex(ctx context.Context, req *proto.DeleteIndexRequest, rsp *proto.DeleteIndexResponse) error {
	_, err := engine.DeleteIndex(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DeleteIndex", err.Error())
	}

	return nil
}

// DeleteAlias Config srv handler
func (cf *Config) DeleteAlias(ctx context.Context, req *proto.DeleteAliasRequest, rsp *proto.DeleteAliasResponse) error {
	_, err := engine.DeleteAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DeleteAlias", err.Error())
	}

	return nil
}

// RenameAlias Config srv handler
func (cf *Config) RenameAlias(ctx context.Context, req *proto.RenameAliasRequest, rsp *proto.RenameAliasResponse) error {
	_, err := engine.RenameAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.RenameAlias", err.Error())
	}

	return nil
}
