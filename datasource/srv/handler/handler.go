package handler

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/validate"
	"github.com/micro/go-micro"
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
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	eng, err := engine.NewDataSourceEngine(req.Endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.NewDataSourceEngine", err.Error())
	}

	datasourcesList, err := engine.SearchDataSources(ctx, srv.Client(), &proto_datasource.SearchRequest{
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
	endpoint, err := eng.Validate(ctx, srv.Client(), datasources)
	if err != nil {
		return errors.BadRequest("go.micro.srv.datasource.eng.Validate", err.Error())
	}

	// Request will be available on After handler wrapper
	// Update req data with the last values
	req.Endpoint = endpoint

	if err := eng.Save(ctx, srv.Client(), endpoint, endpoint.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.eng.Save", err.Error())
	}

	if err := eng.CreateIndexWithAlias(ctx, srv.Client()); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.eng.CreateIndexWithAlias", err.Error())
	}

	if err := s.googleCloudStorage.CreateBucket(endpoint.Index); err != nil {
		return errors.InternalServerError("GoogleCloudStorage", err.Error())
	}

	return nil
}

// Delete datasource handler
func (s *service) Delete(ctx context.Context, req *proto_datasource.DeleteRequest, rsp *proto_datasource.DeleteResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Read datasource
	endpoint, err := engine.ReadDataSource(ctx, srv.Client(), req.Id)
	if err != nil {
		return err
	}

	// Instantiate an engine given datasource
	eng, err := engine.NewDataSourceEngine(endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Request will be available on After handler wrapper
	// Update req data with the last values
	req.Index = endpoint.Index

	// Delete datasource
	if err := eng.Delete(ctx, srv.Client()); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Search datasources handler
func (s *service) Search(ctx context.Context, req *proto_datasource.SearchRequest, rsp *proto_datasource.SearchResponse) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	result, err := engine.SearchDataSources(ctx, srv.Client(), req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	rsp.Result = result.Result
	rsp.Info = result.Info

	return nil
}

// Scan datasources handler, will publish to scan topic to be pick up by crawler srv
func (s *service) Scan(ctx context.Context, req *proto_datasource.ScanRequest, rsp *proto_datasource.ScanResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Read datasource, acts as pre validation before After Handler
	_, err := engine.ReadDataSource(ctx, srv.Client(), req.Id)
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
