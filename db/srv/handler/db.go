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
	_, err := engine.Create(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Create", err.Error())
	}

	return nil
}

// Read db srv handler
func (db *DB) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	response, err := engine.Read(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Read", err.Error())
	}

	rsp.Result = response.Result

	return nil
}

// Update db srv handler
func (db *DB) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	_, err := engine.Update(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Update", err.Error())
	}

	return nil
}

// Delete db srv handler
func (db *DB) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	_, err := engine.Delete(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Delete", err.Error())
	}

	return nil
}

// DeleteByQuery db srv handler
func (db *DB) DeleteByQuery(ctx context.Context, req *proto.DeleteByQueryRequest, rsp *proto.DeleteByQueryResponse) error {
	_, err := engine.DeleteByQuery(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DeleteByQuery", err.Error())
	}

	return nil
}

// Search db srv handler
func (db *DB) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	response, err := engine.Search(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.Search", err.Error())
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}

// Search db srv handler
func (db *DB) SearchById(ctx context.Context, req *proto.SearchByIdRequest, rsp *proto.SearchByIdResponse) error {
	response, err := engine.SearchById(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db.SearchById", err.Error())
	}

	rsp.Result = response.Result

	return nil
}
