package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type EnrichMsgChan struct {
	msg  *enrich_proto.EnrichMessage
	ctx  context.Context
	done chan bool
}

type Enrich struct {
	Client             client.Client
	GoogleCloudStorage *gcslib.GoogleCloudStorage
	EnrichMsgChan      chan EnrichMsgChan
	Workers            int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (e *Enrich) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	c := EnrichMsgChan{
		msg:  enrichmsg,
		ctx:  ctx,
		done: make(chan bool),
	}
	// Queue internally
	e.EnrichMsgChan <- c

	return nil
}

func (e *Enrich) queueListener(wID int) {
	for m := range e.EnrichMsgChan {
		if err := processEnrichMsg(e.Client, e.GoogleCloudStorage, m); err != nil {
			log.Println("Error Processing enrich msg (Audio) on worker", wID, err)
		}
	}
}

func StartWorkers(e *Enrich) {
	for i := 0; i < e.Workers; i++ {
		go e.queueListener(i)
	}
}

func processEnrichMsg(c client.Client, gcs *gcslib.GoogleCloudStorage, m EnrichMsgChan) error {
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

	ch := mfs.AudioEnrich(f, gcs)
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

	m.msg.FileName = f.GetName()

	m.done <- true

	return nil
}
