package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/datasource/srv/engine"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
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
		return errors.InternalServerError("go.micro.srv.datasource GetDataSource", err.Error())
	}

	datasourcesList, err := SearchDataSources(ctx, ds, &proto.SearchRequest{
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
	endpoint, err := eng.Validate(datasources)
	if err != nil {
		return errors.BadRequest("go.micro.srv.datasource", err.Error())
	}

	if err := eng.Save(ctx, endpoint, endpoint.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	if err := CreateIndexWithAlias(ds, ctx, endpoint); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	// Scan created datasource
	if err := ScanDataSource(ds, ctx, endpoint.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Delete datasource handler
func (ds *DataSource) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Id) <= 0 {
		return errors.BadRequest("go.micro.srv.datasource", "id required")
	}

	if err := DeleteDataSource(ctx, ds, req.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.datasource", err.Error())
	}

	return nil
}

// Search datasources handler
func (ds *DataSource) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	result, err := SearchDataSources(ctx, ds, req)
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

// SubscribeCrawlerFinished sets last scan timestamp for the datasource after being scanned and updates crawler state
func SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	var ds *proto.Endpoint

	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, nil)
	rsp, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    msg.DatasourceId,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rsp.Result), &ds); err != nil {
		return err
	}

	ds.CrawlerRunning = false
	ds.LastScan = time.Now().Unix()
	b, err := json.Marshal(ds)
	if err != nil {
		log.Println(err)
	}
	_, err = c.Update(ctx, &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    msg.DatasourceId,
		Data:  string(b),
	})
	if err != nil {
		log.Println(err)
	}

	return nil
}
