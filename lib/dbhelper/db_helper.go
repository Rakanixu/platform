package dbhelper

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

func CreateIntoDB(c client.Client, ctx context.Context, req *db_proto.CreateRequest) (*db_proto.CreateResponse, error) {
	// Inject into context the scope
	nctx := globals.NewDBContext(ctx)

	log.Println("NEWDBCONTECT", nctx)

	creq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Create",
		req,
	)
	crsp := &db_proto.CreateResponse{}
	if err := c.Call(nctx, creq, crsp); err != nil {
		return crsp, err
	}

	return crsp, nil
}

func UpdateFromDB(c client.Client, ctx context.Context, req *db_proto.UpdateRequest) (*db_proto.UpdateResponse, error) {
	dreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Update",
		req,
	)
	drsp := &db_proto.UpdateResponse{}
	if err := c.Call(ctx, dreq, drsp); err != nil {
		return nil, err
	}

	return drsp, nil
}

func DeleteFromDB(c client.Client, ctx context.Context, req *db_proto.DeleteRequest) (*db_proto.DeleteResponse, error) {
	dreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Delete",
		req,
	)
	drsp := &db_proto.DeleteResponse{}
	if err := c.Call(ctx, dreq, drsp); err != nil {
		return nil, err
	}

	return drsp, nil
}

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
