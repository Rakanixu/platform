package engine

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"strings"
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

func cleanFilesHelperIndex(ctx context.Context, endpoint *proto_datasource.Endpoint) error {
	var datasources []*proto_datasource.Endpoint

	// FIXME: pagination
	rsp, err := SearchDataSources(ctx, &proto_datasource.SearchRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
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
		deleteZombieRecords(ctx, datasources, endpoint.Url[:idx])
	}

	return nil
}

func deleteZombieRecords(ctx context.Context, datasources []*proto_datasource.Endpoint, urlToDelete string) {
	delete := 0

	for _, v := range datasources {
		if !strings.Contains(v.Url, urlToDelete) {
			delete++
		}
	}

	if delete >= len(datasources)-1 {
		_, err := operations.Delete(ctx, &proto_operations.DeleteRequest{
			Index: globals.IndexHelper,
			Type:  globals.FileType,
			Id:    globals.GetMD5Hash(urlToDelete[len(localEndpoint):]),
		})
		if err != nil {
			return
		}

		idx := strings.LastIndex(urlToDelete, "/")

		if idx > 0 && urlToDelete[:idx] != "local:/" {
			deleteZombieRecords(ctx, datasources, urlToDelete[:idx])
		}
	}
}
