package handler

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// DataSource struct
type DataSource struct {
	Client             client.Client
	ElasticServiceName string
}

// Create datasource handler
func (ds *DataSource) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Endpoint.Url) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "url required")
	}

	dataSource, err := GetDataSource(ds, req.Endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	if err := dataSource.Validate(); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	if err := dataSource.Save(req.Endpoint, req.Endpoint.Url); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Delete datasource handler
func (ds *DataSource) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Id) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "id required")
	}

	if err := DeleteDataSource(ds, req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Search datasources handler
func (ds *DataSource) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	result, err := SearchDataSources(ds, req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	rsp.Result = result.Result
	rsp.Info = result.Info

	return nil
}

// Scan datasources handler, will publish to scan topic to be pick up by crawler srv
func (ds *DataSource) Scan(ctx context.Context, req *proto.ScanRequest, rsp *proto.ScanResponse) error {
	if len(req.Id) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "id required")
	}

	if err := ScanDataSource(ds, ctx, req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}
