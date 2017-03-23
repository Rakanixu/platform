package subscriber

import (
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	crawler_msg "github.com/kazoup/platform/lib/protomsg/crawler"
	deletebucket_msg "github.com/kazoup/platform/lib/protomsg/deletebucket"
	deletefilebucket_msg "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

type CrawlerFinished struct {
	Client client.Client
	Broker broker.Broker
}

// SubscribeCrawlerFinished sets last scan timestamp for the datasource after being scanned and updates crawler state
func (cf *CrawlerFinished) SubscribeCrawlerFinished(ctx context.Context, msg *crawler_msg.CrawlerFinishedMessage) error {
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
		return err
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

	// Clear index (files that no longer exists, rename, etc..)
	if err := globals.ClearIndex(ctx, cf.Client, ds); err != nil {
		log.Println("ERROR clearing index after scan", err)
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
	Client             client.Client
	Broker             broker.Broker
	GoogleCloudStorage *gcslib.GoogleCloudStorage
}

// SubscribeDeleteBucket subscribes to DeleteBucket Message to clean un a bicket in GC storage
func (db *DeleteBucket) SubscribeDeleteBucket(ctx context.Context, msg *deletebucket_msg.DeleteBucketMsg) error {
	return db.GoogleCloudStorage.DeleteBucket(msg.Index)
}

type DeleteFileInBucket struct {
	Client             client.Client
	Broker             broker.Broker
	GoogleCloudStorage *gcslib.GoogleCloudStorage
}

// SubscribeCleanBucket subscribes to DCleanBucket Message to remove thumbs not longer related with document in index
func (dfb *DeleteFileInBucket) SubscribeDeleteFileInBucket(ctx context.Context, msg *deletefilebucket_msg.DeleteFileInBucketMsg) error {
	return dfb.GoogleCloudStorage.Delete(msg.Index, msg.FileId)
}
