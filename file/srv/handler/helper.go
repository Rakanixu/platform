package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
)

func NewFileSystem(ctx context.Context, id string) (fs.Fs, error) {
	var ds *proto_datasource.Endpoint

	// Get the datasource
	rsp, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    id,
	})

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(rsp.Result), &ds); err != nil {
		return nil, err
	}

	// Instantiate file system from datasource
	fsys, err := fs.NewFsFromEndpoint(ds)
	if err != nil {
		return nil, err
	}

	return fsys, nil
}
