package handler

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func UpdateFileSystemAuth(fc client.Client, ctx context.Context, id string, token *datasource_proto.Token) error {
	var ds *datasource_proto.Endpoint

	// Get the datasource
	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, fc)
	rr, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    id,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rr.Result), &ds); err != nil {
		return err
	}

	ds.Token = token

	b, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	_, err = c.Update(ctx, &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}
