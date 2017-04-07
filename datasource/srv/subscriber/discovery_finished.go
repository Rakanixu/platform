package subscriber

import (
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	crawler "github.com/kazoup/platform/lib/protomsg/crawler"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"time"
)

type DiscoveryFinished struct{}

// PostDiscovery sets last scan timestamp for the datasource after being scanned and updates crawler state
func (df *DiscoveryFinished) PostDiscovery(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, srv.Client())
	rsp, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    msg.DatasourceId,
	})
	if err != nil {
		return err
	}

	var ds *proto.Endpoint
	if err := json.Unmarshal([]byte(rsp.Result), &ds); err != nil {
		return err
	}

	ds.CrawlerRunning = false
	ds.LastScan = time.Now().Unix()
	b, err := json.Marshal(ds)
	if err != nil {
		return err
	}
	_, err = c.Update(ctx, &db_proto.UpdateRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    msg.DatasourceId,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	// Clear index (files that no longer exists, rename, etc..)
	if err := globals.ClearIndex(ctx, srv.Client(), ds); err != nil {
		return err
	}

	return nil
}
