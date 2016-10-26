package handler

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/structs/fs"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// File struct
type File struct {
	Client client.Client
}

// Create File handler
func (f *File) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	var ds *datasource_proto.Endpoint

	if len(req.DatasourceId) == 0 {
		return errors.BadRequest("com.kazoup.srv.file", "datasource_id required")
	}

	// Get the datasource
	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, f.Client)
	rr, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    req.DatasourceId,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rr.Result), &ds); err != nil {
		return err
	}

	// Instantiate file system from datasource
	fsys, err := fs.NewFsFromEndpoint(ds)
	if err != nil {
		return err
	}

	// Create a file for given file system
	s, err := fsys.CreateFile(req.MimeType)
	if err != nil {
		return err
	}

	rsp.DocUrl = s

	return nil
}
