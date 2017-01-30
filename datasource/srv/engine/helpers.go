package engine

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"strings"
)

// ReadDataSource returns the endpoint with given id
func ReadDataSource(ctx context.Context, c client.Client, id string) (*datasource_proto.Endpoint, error) {
	srvReq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	srvRes := &db_proto.ReadResponse{}

	if err := c.Call(ctx, srvReq, srvRes); err != nil {
		return nil, err
	}

	var endpoint *datasource_proto.Endpoint

	if err := json.Unmarshal([]byte(srvRes.Result), &endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// SearchDataSources queries for datasources stored in ES
func SearchDataSources(ctx context.Context, c client.Client, req *datasource_proto.SearchRequest) (*datasource_proto.SearchResponse, error) {
	srvReq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Search",
		&db_proto.SearchRequest{
			Index:    "datasources",
			Type:     "datasource",
			From:     req.From,
			Size:     req.Size,
			Category: req.Category,
			Term:     req.Term,
			Depth:    req.Depth,
			Url:      req.Url,
		},
	)
	srvRes := &db_proto.SearchResponse{}

	if err := c.Call(ctx, srvReq, srvRes); err != nil {
		return nil, err
	}

	rsp := &datasource_proto.SearchResponse{
		Result: srvRes.Result,
		Info:   srvRes.Info,
	}

	return rsp, nil
}

func cleanFilesHelperIndex(ctx context.Context, c client.Client, endpoint *datasource_proto.Endpoint) error {
	var datasources []*datasource_proto.Endpoint

	rsp, err := SearchDataSources(ctx, c, &datasource_proto.SearchRequest{
		Index: "datasources",
		Type:  "datasource",
		From:  0,
		Size:  9999,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rsp.Result), &datasources); err != nil {
		return err
	}

	idx := strings.LastIndex(endpoint.Url, "/")
	if idx > 0 {
		deleteZombieRecords(ctx, c, datasources, endpoint.Url[:idx])
	}

	return nil
}

func deleteZombieRecords(ctx context.Context, c client.Client, datasources []*datasource_proto.Endpoint, urlToDelete string) {
	delete := 0

	for _, v := range datasources {
		if !strings.Contains(v.Url, urlToDelete) {
			delete++
		}
	}

	if delete >= len(datasources)-1 {
		deleteReq := c.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: globals.IndexHelper,
				Type:  "file",
				Id:    globals.GetMD5Hash(urlToDelete[len(localEndpoint):]),
			},
		)
		deleteRes := &db_proto.DeleteResponse{}

		if err := c.Call(ctx, deleteReq, deleteRes); err != nil {
			log.Println("ERROR", err)
		}
		idx := strings.LastIndex(urlToDelete, "/")

		if idx > 0 && urlToDelete[:idx] != "local:/" {
			deleteZombieRecords(ctx, c, datasources, urlToDelete[:idx])
		}
	}
}
