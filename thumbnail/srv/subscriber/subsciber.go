package subscriber

import (
	"encoding/json"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

func NewTaskHandler(workers int, gcs *gcslib.GoogleCloudStorage) *taskHandler {
	t := &taskHandler{
		googleCloudStorage: gcs,
		thumbnailMsgChan:   make(chan thumbnailMsgChan, 1000000),
		Workers:            workers,
	}

	startWorkers(t)

	return t
}

type taskHandler struct {
	googleCloudStorage *gcslib.GoogleCloudStorage
	thumbnailMsgChan   chan thumbnailMsgChan
	Workers            int
}

type thumbnailMsgChan struct {
	ctx context.Context
	msg *enrich.EnrichMessage
	err chan error
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (t *taskHandler) Thumbnail(ctx context.Context, enrichmsg *enrich.EnrichMessage) error {
	c := thumbnailMsgChan{
		ctx: ctx,
		msg: enrichmsg,
		err: make(chan error),
	}
	// Queue internally
	t.thumbnailMsgChan <- c

	return <-c.err
}

// queueListener range over thumbnailMsgChan channel and process msgs one by one
func (t *taskHandler) queueListener(wID int) {
	for m := range t.thumbnailMsgChan {
		if err := processThumbnailMsg(t.googleCloudStorage, m); err != nil {
			m.err <- errors.NewPlatformError(globals.THUMBNAIL_SERVICE_NAME, "processEnrichMsg", fmt.Sprintf("worker %d", wID), err)
		}
		// Successful
		m.err <- nil
	}
}

func startWorkers(t *taskHandler) {
	// Start workers
	for i := 0; i < t.Workers; i++ {
		go t.queueListener(i)
	}
}

func processThumbnailMsg(gcs *gcslib.GoogleCloudStorage, m thumbnailMsgChan) error {
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

	ch := mfs.Thumbnail(f, gcs)
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

	return nil
}
