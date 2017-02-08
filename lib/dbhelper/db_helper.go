package dbhelper

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func CreateIntoDB(c client.Client, ctx context.Context, req *db_proto.CreateRequest) (*db_proto.CreateResponse, error) {
	creq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Create",
		req,
	)
	crsp := &db_proto.CreateResponse{}
	if err := c.Call(ctx, creq, crsp); err != nil {
		return crsp, err
	}

	return crsp, nil
}

func ReadFromDB(c client.Client, ctx context.Context, req *db_proto.ReadRequest) (*db_proto.ReadResponse, error) {
	rreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Read",
		req,
	)
	rrsp := &db_proto.ReadResponse{}
	if err := c.Call(ctx, rreq, rrsp); err != nil {
		return rrsp, err
	}

	return rrsp, nil
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

func SearchFromDB(c client.Client, ctx context.Context, req *db_proto.SearchRequest) (*db_proto.SearchResponse, error) {
	sreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Search",
		req,
	)
	srsp := &db_proto.SearchResponse{}
	if err := c.Call(ctx, sreq, srsp); err != nil {
		return nil, err
	}

	return srsp, nil
}

func SearchById(c client.Client, ctx context.Context, req *db_proto.SearchByIdRequest) (*db_proto.SearchByIdResponse, error) {
	sreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.SearchById",
		req,
	)
	srsp := &db_proto.SearchByIdResponse{}
	if err := c.Call(ctx, sreq, srsp); err != nil {
		return nil, err
	}

	return srsp, nil
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
