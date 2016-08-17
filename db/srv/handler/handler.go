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
	rsp, err := engine.Create(req)
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
