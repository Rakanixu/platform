package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func NewFileSystem(fc client.Client, ctx context.Context, id string) (fs.Fs, error) {
	var ds *proto_datasource.Endpoint

	// Get the datasource
	rsp, err := db_helper.ReadFromDB(fc, ctx, &db_proto.ReadRequest{
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
