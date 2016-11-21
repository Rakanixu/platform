package handler

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	scheduler_proto "github.com/kazoup/platform/scheduler/srv/proto/scheduler"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"time"
)

// DataSource struct
type DataSource struct {
	Client client.Client
}

// Create datasource handler
func (ds *DataSource) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Endpoint.Url) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "url required")
	}

	eng, err := engine.NewDataSourceEngine(req.Endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	datasourcesList, err := engine.SearchDataSources(ctx, ds.Client, &proto.SearchRequest{
		Index: "datasources",
		Type:  "datasource",
		From:  0,
		Size:  9999,
	})

	datasources := "[]"
	if datasourcesList != nil {
		datasources = datasourcesList.Result
	}

	// Validate and assigns Id and index
	endpoint, err := eng.Validate(ctx, ds.Client, datasources)
	if err != nil {
		return errors.BadRequest("go.micro.srv.datasource", err.Error())
	}

	if err := eng.Save(ctx, endpoint, endpoint.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	if err := eng.CreateIndexWithAlias(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Scan created datasource
	if err := eng.Scan(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Shedule scan task
	if err := eng.ScheduleScan(ctx, ds.Client, &scheduler_proto.CreateScheduledTaskRequest{
		Task: &scheduler_proto.Task{
			Id:     endpoint.Id,
			Action: globals.StartScanTask,
		},
		Schedule: &scheduler_proto.Schedule{
			IntervalSeconds: int64(time.Hour.Seconds()),
		},
	}); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Delete datasource handler
func (ds *DataSource) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Id) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "id required")
	}

	// Read datasource
	endpoint, err := engine.ReadDataSource(ctx, ds.Client, req.Id)
	if err != nil {
		return err
	}

	// Instantiate an engine given datasource
	eng, err := engine.NewDataSourceEngine(endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Delete datasource
	if err := eng.Delete(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Search datasources handler
func (ds *DataSource) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	result, err := engine.SearchDataSources(ctx, ds.Client, req)
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

	// Read datasource
	endpoint, err := engine.ReadDataSource(ctx, ds.Client, req.Id)
	if err != nil {
		return err
	}

	// Instantiate an engine given datasource
	eng, err := engine.NewDataSourceEngine(endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Start scan
	if err := eng.Scan(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}
