package handler

import (
	"github.com/kazoup/platform/datasource/srv/engine"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

// DataSource struct
type DataSource struct {
	Client             client.Client
	GoogleCloudStorage *gcslib.GoogleCloudStorage
}

// Create datasource handler
func (ds *DataSource) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Endpoint.Url) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "url required")
	}

	eng, err := engine.NewDataSourceEngine(req.Endpoint)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.NewDataSourceEngine", err.Error())
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
		return errors.BadRequest("go.micro.srv.datasource.eng.Validate", err.Error())
	}

	if err := eng.Save(ctx, ds.Client, endpoint, endpoint.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.eng.Save", err.Error())
	}

	if err := eng.CreateIndexWithAlias(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.eng.CreateIndexWithAlias", err.Error())
	}

	if err := ds.GoogleCloudStorage.CreateBucket(endpoint.Index); err != nil {
		return errors.InternalServerError("GoogleCloudStorage", err.Error())
	}

	// Scan created datasource
	if err := eng.Scan(ctx, ds.Client); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource.eng.Scan", err.Error())
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

	// Publish message to clean async the bucket that stores the thumbnails in GC storage
	if err := ds.Client.Publish(globals.NewSystemContext(), ds.Client.NewPublication(globals.DeleteBucketTopic, &proto.DeleteBucketMessage{
		Endpoint: endpoint,
	})); err != nil {
		log.Println("ERROR cleaningthumbs from GCS", err)
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

// ScanAll datasources handler, will publish to scan topic
// If req.DatasourcesId not empty, those specific datasources will be scanned
// If req.DatasourcesId empty, all user datasources will be scanned
func (ds *DataSource) ScanAll(ctx context.Context, req *proto.ScanAllRequest, rsp *proto.ScanAllResponse) error {
	if len(req.DatasourcesId) > 0 {
		// Scan all datasources specified on request
		for _, v := range req.DatasourcesId {
			if err := ds.Scan(ctx, &proto.ScanRequest{Id: v}, &proto.ScanResponse{}); err != nil {
				log.Println("ERROR starting scan for ", v, err)
			}
		}
	} else {
		// Scan all datasources for given user
		uID, err := globals.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		if err := engine.ScanAllDatasources(ctx, ds.Client, uID); err != nil {
			return errors.InternalServerError("go.micro.srv.datasource.ScanAll", err.Error())
		}
	}

	return nil
}

func (ds *DataSource) Health(ctx context.Context, req *proto.HealthRequest, rsp *proto.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
