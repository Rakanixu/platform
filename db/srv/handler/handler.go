package handler

import (
	engine "github.com/kazoup/platform/db/srv/engine"
	proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// DB struct
type DB struct{}

// Create db srv handler
func (db *DB) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	_, err := engine.Create(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Read db srv handler
func (db *DB) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	response, err := engine.Read(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Result = response.Result

	return nil
}

// Update db srv handler
func (db *DB) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	_, err := engine.Update(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Delete db srv handler
func (db *DB) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	_, err := engine.Delete(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// CreateIndexWithSettings db srv handler
func (db *DB) CreateIndexWithSettings(ctx context.Context, req *proto.CreateIndexWithSettingsRequest, rsp *proto.CreateIndexWithSettingsResponse) error {
	_, err := engine.CreateIndexWithSettings(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// CreateIndexWithSettings db srv handler
func (db *DB) PutMappingFromJSON(ctx context.Context, req *proto.PutMappingFromJSONRequest, rsp *proto.PutMappingFromJSONResponse) error {
	_, err := engine.PutMappingFromJSON(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Status db srv handler
func (db *DB) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	response, err := engine.Status(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Status = response.Status

	return nil
}
