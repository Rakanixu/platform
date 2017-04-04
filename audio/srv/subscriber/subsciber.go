package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

func NewTaskHandler(workers int, cs *gcslib.GoogleCloudStorage) *taskHandler {
	t := &taskHandler{
		googleCloudStorage: cs,
		enrichMsgChan:      make(chan enrichMsgChan, 1000000),
		workers:            workers,
	}

	startWorkers(t)

	return t
}

type taskHandler struct {
	googleCloudStorage *gcslib.GoogleCloudStorage
	enrichMsgChan      chan enrichMsgChan
	workers            int
}

type enrichMsgChan struct {
	msg  *enrich_proto.EnrichMessage
	ctx  context.Context
	done chan bool
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (e *taskHandler) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	c := enrichMsgChan{
		msg:  enrichmsg,
		ctx:  ctx,
		done: make(chan bool),
	}
	// Queue internally
	e.enrichMsgChan <- c

	<-c.done

	return nil
}

func (e *taskHandler) queueListener(wID int) {
	for m := range e.enrichMsgChan {
		if err := processEnrichMsg(e.googleCloudStorage, m); err != nil {
			log.Println("Error Processing enrich msg (Audio) on worker", wID, err)
		}
	}
}

func startWorkers(t *taskHandler) {
	for i := 0; i < t.workers; i++ {
		go t.queueListener(i)
	}
}

func processEnrichMsg(gcs *gcslib.GoogleCloudStorage, m enrichMsgChan) error {
	srv, ok := micro.FromContext(m.ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

	frsp, err := db_helper.ReadFromDB(srv.Client(), m.ctx, &db_proto.ReadRequest{
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

	drsp, err := db_helper.ReadFromDB(srv.Client(), m.ctx, &db_proto.ReadRequest{
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

	_, err = db_helper.UpdateFromDB(srv.Client(), m.ctx, &db_proto.UpdateRequest{
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
