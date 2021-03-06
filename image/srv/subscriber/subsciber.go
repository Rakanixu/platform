package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"golang.org/x/net/context"
)

func NewTaskHandler(workers int) *taskHandler {
	t := &taskHandler{
		enrichMsgChan: make(chan enrichMsgChan, 1000000),
		workers:       workers,
	}

	startWorkers(t)

	return t
}

type taskHandler struct {
	enrichMsgChan chan enrichMsgChan
	workers       int
}

type enrichMsgChan struct {
	ctx context.Context
	msg *enrich.EnrichMessage
	err chan error
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (t *taskHandler) Enrich(ctx context.Context, enrichmsg *enrich.EnrichMessage) error {
	c := enrichMsgChan{
		ctx: ctx,
		msg: enrichmsg,
		err: make(chan error),
	}

	// Queue internally
	t.enrichMsgChan <- c

	return <-c.err
}

// queueListener range over enrichMsgChan channel and process msgs one by one
func (t *taskHandler) queueListener(wID int) {
	for m := range t.enrichMsgChan {
		if err := processEnrichMsg(m); err != nil {
			m.err <- errors.NewPlatformError(globals.IMAGE_SERVICE_NAME, "processEnrichMsg", fmt.Sprintf("worker %d", wID), 500, err)
		}
		// Successful
		m.err <- nil
	}
}

func startWorkers(t *taskHandler) {
	// Start workers
	for i := 0; i < t.workers; i++ {
		go t.queueListener(i)
	}
}

func processEnrichMsg(m enrichMsgChan) error {
	frsp, err := operations.Read(m.ctx, &proto_operations.ReadRequest{
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

	drsp, err := operations.Read(m.ctx, &proto_operations.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    f.GetDatasourceID(),
	})
	if err != nil {
		return err
	}

	var endpoint proto_datasource.Endpoint
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

	_, err = operations.Update(m.ctx, &proto_operations.UpdateRequest{
		Index: m.msg.Index,
		Type:  globals.FileType,
		Id:    m.msg.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	m.msg.FileName = f.GetName()

	return nil
}
