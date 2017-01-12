package handler

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func NewFileSystem(fc client.Client, ctx context.Context, id string) (fs.Fs, error) {
	var ds *datasource_proto.Endpoint

	// Get the datasource
	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, fc)
	rr, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    id,
	})
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(rr.Result), &ds); err != nil {
		return nil, err
	}

	// Instantiate file system from datasource
	fsys, err := fs.NewFsFromEndpoint(ds)
	if err != nil {
		return nil, err
	}

	return fsys, nil
}
