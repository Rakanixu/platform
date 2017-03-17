package subscriber

import (
	"encoding/json"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type EnrichMsgChan struct {
	ctx context.Context
	msg *enrich_proto.EnrichMessage
}

type Enrich struct {
	Client        client.Client
	EnrichMsgChan chan EnrichMsgChan
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (e *Enrich) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	// Queue internally
	e.EnrichMsgChan <- EnrichMsgChan{
		ctx: ctx,
		msg: enrichmsg,
	}

	return nil
}

func (e *Enrich) queueListener(wID int) {
	for m := range e.EnrichMsgChan {
		if err := processEnrichMsg(e.Client, m); err != nil {
			log.Println("Error Processing enrich msg (Document) on Worker", wID, err)
		}
	}
}

func StartWorkers(e *Enrich) {
	for i := 0; i < e.Workers; i++ {
		go e.queueListener(i)
	}
}

func processEnrichMsg(c client.Client, m EnrichMsgChan) error {
	frsp, err := db_helper.ReadFromDB(c, m.ctx, &db_proto.ReadRequest{
		Index: m.msg.Index,
		Type:  globals.FileType,
		Id:    m.msg.Id,
	})
	if err != nil {
		return err
	}

	f, err := file.NewFileFromString(frsp.Result)
	if err != nil {
		return err
	}

	drsp, err := db_helper.ReadFromDB(c, m.ctx, &db_proto.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    f.GetDatasourceID(),
	})
	if err != nil {
		return err
	}

	var endpoint datasource_proto.Endpoint
	if err := json.Unmarshal([]byte(drsp.Result), &endpoint); err != nil {
		return err
	}

	mfs, err := fs.NewFsFromEndpoint(&endpoint)
	if err != nil {
		return err
	}

	ch := mfs.DocEnrich(f)
	// Block while enriching, we expect only one m
	fm := <-ch
	close(ch)

	if fm.Error != err {
		return fm.Error
	}

	b, err := json.Marshal(fm.File)
	if err != nil {
		return err
	}

	_, err = db_helper.UpdateFromDB(c, m.ctx, &db_proto.UpdateRequest{
		Index: m.msg.Index,
		Type:  globals.FileType,
		Id:    m.msg.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	// Publish notification topic if requested
	if m.msg.Notify {
		if err := c.Publish(m.ctx, c.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
			Method: globals.NOTIFY_REFRESH_SEARCH,
			UserId: m.msg.UserId,
			Info:   fmt.Sprintf("Document content extraction for %s finished.", f.GetName()),
		})); err != nil {
			log.Print("Publishing NotificationTopic (DocEnrich) error %s", err)
		}
	}

	// Publish the same message to ExtractEntitiesTopic
	if err := c.Publish(m.ctx, c.NewPublication(globals.ExtractEntitiesTopic, m.msg)); err != nil {
		return err
	}

	// Publish the same message to SentimentEnrichTopic
	if err := c.Publish(m.ctx, c.NewPublication(globals.SentimentEnrichTopic, m.msg)); err != nil {
		return err
	}

	return nil
}
