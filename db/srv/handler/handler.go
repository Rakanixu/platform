package handler

import (
	engine "github.com/kazoup/platform/db/srv/engine"
	proto "github.com/kazoup/platform/db/srv/proto/db"
	search_proto "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// DB struct
type DB struct{}

// Create db srv handler
func (db *DB) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	_, err := engine.Create(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Read db srv handler
func (db *DB) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	response, err := engine.Read(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Result = response.Result

	return nil
}

// Update db srv handler
func (db *DB) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	_, err := engine.Update(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Delete db srv handler
func (db *DB) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	_, err := engine.Delete(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// DeleteByQuery db srv handler
func (db *DB) DeleteByQuery(ctx context.Context, req *proto.DeleteByQueryRequest, rsp *proto.DeleteByQueryResponse) error {
	_, err := engine.DeleteByQuery(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// CreateIndexWithSettings db srv handler
func (db *DB) CreateIndexWithSettings(ctx context.Context, req *proto.CreateIndexWithSettingsRequest, rsp *proto.CreateIndexWithSettingsResponse) error {
	_, err := engine.CreateIndexWithSettings(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// CreateIndexWithSettings db srv handler
func (db *DB) PutMappingFromJSON(ctx context.Context, req *proto.PutMappingFromJSONRequest, rsp *proto.PutMappingFromJSONResponse) error {
	_, err := engine.PutMappingFromJSON(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Status db srv handler
func (db *DB) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	response, err := engine.Status(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Status = response.Status

	return nil
}

// Search db srv handler
func (db *DB) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	response, err := engine.Search(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}

// Search db srv handler
func (db *DB) SearchById(ctx context.Context, req *proto.SearchByIdRequest, rsp *proto.SearchByIdResponse) error {
	response, err := engine.SearchById(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Result = response.Result

	return nil
}

// RenameIndexAlias db srv handler
func (db *DB) AddAlias(ctx context.Context, req *proto.AddAliasRequest, rsp *proto.AddAliasResponse) error {
	_, err := engine.AddAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// DeleteIndex db srv handler
func (db *DB) DeleteIndex(ctx context.Context, req *proto.DeleteIndexRequest, rsp *proto.DeleteIndexResponse) error {
	_, err := engine.DeleteIndex(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// DeleteAlias db srv handler
func (db *DB) DeleteAlias(ctx context.Context, req *proto.DeleteAliasRequest, rsp *proto.DeleteAliasResponse) error {
	_, err := engine.DeleteAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// RenameAlias db srv handler
func (db *DB) RenameAlias(ctx context.Context, req *proto.RenameAliasRequest, rsp *proto.RenameAliasResponse) error {
	_, err := engine.RenameAlias(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	return nil
}

// Aggregate db srv handler
func (db *DB) Aggregate(ctx context.Context, req *search_proto.AggregateRequest, rsp *search_proto.AggregateResponse) error {
	response, err := engine.Aggregate(ctx, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.db", err.Error())
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}
