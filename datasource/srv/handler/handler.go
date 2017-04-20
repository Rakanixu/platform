package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/engine"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/validate"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

func NewServiceHandler(cloudStorage *gcslib.GoogleCloudStorage) *service {
	return &service{
		googleCloudStorage: cloudStorage,
	}
}

type service struct {
	googleCloudStorage *gcslib.GoogleCloudStorage
}

// Create datasource handler
func (s *service) Create(ctx context.Context, req *proto_datasource.CreateRequest, rsp *proto_datasource.CreateResponse) error {
	if err := validate.Exists(ctx, req.Endpoint.Url); err != nil {
		return errors.BadRequest(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	eng, err := engine.NewDataSourceEngine(req.Endpoint)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	datasourcesList, err := engine.SearchDataSources(ctx, &proto_datasource.SearchRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		From:  0,
		Size:  9999,
	})

	datasources := "[]"
	if datasourcesList != nil {
		datasources = datasourcesList.Result
	}

	// Validate and assigns Id and index
	endpoint, err := eng.Validate(ctx, datasources)
	if err != nil {
		return errors.BadRequest(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Request will be available on After handler wrapper
	// Update req data with the last values
	req.Endpoint = endpoint

	if err := eng.Save(ctx, endpoint, endpoint.Id); err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	if err := eng.CreateIndexWithAlias(ctx); err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	if err := s.googleCloudStorage.CreateBucket(endpoint.Index); err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	return nil
}

// Read datasource handler
func (s *service) Read(ctx context.Context, req *proto_datasource.ReadRequest, rsp *proto_datasource.ReadResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return errors.BadRequest(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Read datasource
	endpoint, err := engine.ReadDataSource(ctx, req.Id)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	b, err := json.Marshal(endpoint)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	rsp.Result = string(b)

	return nil
}

// Delete datasource handler
func (s *service) Delete(ctx context.Context, req *proto_datasource.DeleteRequest, rsp *proto_datasource.DeleteResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return errors.BadRequest(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Read datasource
	endpoint, err := engine.ReadDataSource(ctx, req.Id)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Instantiate an engine given datasource
	eng, err := engine.NewDataSourceEngine(endpoint)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Request will be available on After handler wrapper
	// Update req data with the last values
	req.Index = endpoint.Index

	// Delete datasource
	if err := eng.Delete(ctx); err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	return nil
}

// Search datasources handler
func (s *service) Search(ctx context.Context, req *proto_datasource.SearchRequest, rsp *proto_datasource.SearchResponse) error {
	result, err := engine.SearchDataSources(ctx, req)
	if err != nil {
		return errors.InternalServerError(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	rsp.Result = result.Result
	rsp.Info = result.Info

	return nil
}

// Scan datasources handler, will publish to scan topic to be pick up by crawler srv
func (s *service) Scan(ctx context.Context, req *proto_datasource.ScanRequest, rsp *proto_datasource.ScanResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return errors.BadRequest(globals.DATASOURCE_SERVICE_NAME, err.Error())
	}

	// Read datasource, acts as pre validation before After Handler
	_, err := engine.ReadDataSource(ctx, req.Id)
	if err != nil {
		return err
	}

	return nil
}

// ScanAll datasources handler
// If req.DatasourcesId not empty, those specific datasources will be scanned
// If req.DatasourcesId empty, all user datasources will be scanned
func (ds *service) ScanAll(ctx context.Context, req *proto_datasource.ScanAllRequest, rsp *proto_datasource.ScanAllResponse) error {
	// Aknowledge

	return nil
}

func (ds *service) Health(ctx context.Context, req *proto_datasource.HealthRequest, rsp *proto_datasource.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
