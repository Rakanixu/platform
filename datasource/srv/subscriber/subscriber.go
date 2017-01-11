package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

type CrawlerStarted struct {
	Client client.Client
	Broker broker.Broker
}

// SubscribeCrawlerStarted receives CrawlerStartedMessage and publish to NotificationTopic
func (cs *CrawlerStarted) SubscribeCrawlerStarted(ctx context.Context, msg *crawler.CrawlerStartedMessage) error {
	var ds *proto.Endpoint

	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, cs.Client)
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

	// Publish notification
	nm := &notification_proto.NotificationMessage{
		Info:   "Scan started on " + ds.Url + " datasource.",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: msg.UserId,
		Data:   rsp.Result,
	}

	// Publish notification
	if err := cs.Client.Publish(ctx, cs.Client.NewPublication(globals.NotificationTopic, nm)); err != nil {
		return err
	}

	return nil
}

type CrawlerFinished struct {
	Client client.Client
	Broker broker.Broker
}

// SubscribeCrawlerFinished sets last scan timestamp for the datasource after being scanned and updates crawler state
func (cf *CrawlerFinished) SubscribeCrawlerFinished(ctx context.Context, msg *crawler.CrawlerFinishedMessage) error {
	var ds *proto.Endpoint

	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, cf.Client)
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
	if err := cf.Client.Publish(ctx, cf.Client.NewPublication(globals.NotificationTopic, nm)); err != nil {
		return err
	}

	return nil
}

type DeleteBucket struct {
	Client client.Client
	Broker broker.Broker
}

// SubscribeDeleteBucket subscribes to DeleteBucket Message to clean un a bicket in GC storage
func (db *DeleteBucket) SubscribeDeleteBucket(ctx context.Context, msg *proto.DeleteBucketMessage) error {
	ncs, err := cs.NewCloudStorageFromEndpoint(msg.Endpoint, globals.GoogleCloudStorage)
	if err != nil {
		return err
	}

	return ncs.DeleteBucket()
}

type DeleteFileInBucket struct {
	Client client.Client
	Broker broker.Broker
}

// SubscribeCleanBucket subscribes to DCleanBucket Message to remove thumbs not longer related with document in index
func (dfb *DeleteFileInBucket) SubscribeDeleteFileInBucket(ctx context.Context, msg *proto.DeleteFileInBucketMessage) error {
	return fs.DeleteFile(msg.Index, msg.FileId)
}
