package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type Enrich struct {
	Client        client.Client
	EnrichMsgChan chan *enrich_proto.EnrichMessage
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (e *Enrich) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	// Queue internally
	e.EnrichMsgChan <- enrichmsg

	return nil
}

// queueListener range over EnrichMsgChan channel and process msgs one by one
func (e *Enrich) queueListener(wID int) {
	for m := range e.EnrichMsgChan {
		if err := processEnrichMsg(e.Client, m); err != nil {
			log.Println("Error Processing enrich msg (Image) on worker ", wID, err)
		}
	}
}

func StartWorkers(e *Enrich) {
	// Start workers
	for i := 0; i < e.Workers; i++ {
		go e.queueListener(i)
	}
}

func processEnrichMsg(c client.Client, m *enrich_proto.EnrichMessage) error {
	frsp, err := db_helper.ReadFromDB(c, globals.NewSystemContext(), &db_proto.ReadRequest{
		Index: m.Index,
		Type:  globals.FileType,
		Id:    m.Id,
	})
	if err != nil {
		return err
	}

	f, err := file.NewFileFromString(frsp.Result)
	if err != nil {
		return err
	}

	drsp, err := db_helper.ReadFromDB(c, globals.NewSystemContext(), &db_proto.ReadRequest{
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

	_, err = db_helper.UpdateFromDB(c, globals.NewSystemContext(), &db_proto.UpdateRequest{
		Index: m.Index,
		Type:  globals.FileType,
		Id:    m.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}
