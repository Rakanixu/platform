package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

var Broker broker.Broker

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

	// Publish notification
	nm := &notification_proto.NotificationMessage{
		Info:   "Scan finished on " + ds.Url + " datasource.",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: ds.UserId,
		Data:   string(b),
	}

	// Publish notification
	if err := client.Publish(ctx, client.NewPublication(globals.NotificationTopic, nm)); err != nil {
		return err
	}

	return nil
}
