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
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type EnrichMsgChan struct {
	ctx  context.Context
	msg  *enrich_proto.EnrichMessage
	done chan bool
}

type TaskHandler struct {
	EnrichMsgChan chan EnrichMsgChan
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (t *TaskHandler) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	c := EnrichMsgChan{
		ctx:  ctx,
		msg:  enrichmsg,
		done: make(chan bool),
	}

	// Queue internally
	t.EnrichMsgChan <- c

	<-c.done

	return nil
}

// queueListener range over EnrichMsgChan channel and process msgs one by one
func (t *TaskHandler) queueListener(wID int) {
	for m := range t.EnrichMsgChan {
		if err := processEnrichMsg(m); err != nil {
			log.Println("Error Processing enrich msg (Image) on worker ", wID, err)
		}
	}
}

func StartWorkers(t *TaskHandler) {
	// Start workers
	for i := 0; i < t.Workers; i++ {
		go t.queueListener(i)
	}
}

func processEnrichMsg(m EnrichMsgChan) error {
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

	ch := mfs.ImgEnrich(f)
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
