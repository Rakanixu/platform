package engine

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
)

// ReadDataSource returns the endpoint with given id
func ReadDataSource(ctx context.Context, id string) (*proto_datasource.Endpoint, error) {
	rsp, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    id,
	})
	if err != nil {
		return nil, err
	}

	var endpoint *proto_datasource.Endpoint

	if err := json.Unmarshal([]byte(rsp.Result), &endpoint); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// SearchDataSources queries for datasources stored in ES
func SearchDataSources(ctx context.Context, req *proto_datasource.SearchRequest) (*proto_datasource.SearchResponse, error) {
	rsp, err := operations.Search(ctx, &proto_operations.SearchRequest{
		Index:    globals.IndexDatasources,
		Type:     globals.TypeDatasource,
		From:     req.From,
		Size:     req.Size,
		Category: req.Category,
		Term:     req.Term,
		Depth:    req.Depth,
		Url:      req.Url,
	})
	if err != nil {
		return nil, err
	}

	ds := &proto_datasource.SearchResponse{
		Result: rsp.Result,
		Info:   rsp.Info,
	}

	return ds, nil
}
