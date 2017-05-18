package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/normalization/text"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/lib/rossete"
	"golang.org/x/net/context"
	"strings"
	"time"
)

const (
	//https://developer.rosette.com/features-and-functions#-entity-types
	IDENTIFIER = "IDENTIFIER"
	PERSON     = "PERSON"
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

func (t *taskHandler) queueListener(wID int) {
	for m := range t.enrichMsgChan {
		if err := processEnrichMsg(m); err != nil {
			m.err <- errors.NewPlatformError(globals.ENTITIES_SERVICE_NAME, "processEnrichMsg", fmt.Sprintf("worker %d", wID), 500, err)
		}
		// Successful
		m.err <- nil
	}
}

func startWorkers(t *taskHandler) {
	for i := 0; i < t.workers; i++ {
		go t.queueListener(i)
	}
}

func processEnrichMsg(m enrichMsgChan) error {
	rsp, err := operations.Read(m.ctx, &proto_operations.ReadRequest{
		Index: m.msg.Index,
		Type:  globals.FileType,
		Id:    m.msg.Id,
	})
	if err != nil {
		return err
	}

	f, err := file.NewFileFromString(rsp.Result)
	if err != nil {
		return err
	}

	process := false
	tm := f.GetOptsTimestamps()

	if tm == nil || tm.TextAnalyzedTimestamp == nil {
		process = true
	} else {
		process = tm.TextAnalyzedTimestamp.Before(f.GetModifiedTime())
	}

	if process {
		// Apply rossete
		if len(f.GetContent()) > 0 {
			// Sanitaze content
			// We do not save sanitazed because if we want to display, it maintains some format
			t, err := text.ReplaceDoubleQuotes(f.GetContent())
			if err != nil {
				return err
			}
			t, err = text.ReplaceTabs(t)
			if err != nil {
				return err
			}
			t, err = text.ReplaceNewLines(t)
			if err != nil {
				return err
			}

			e, err := rossete.Entities(t)
			if err != nil {
				return err
			}
			f.SetEntities(e)

			for _, ent := range e.Entities {
				if strings.Contains(ent.Type, IDENTIFIER) || strings.Contains(ent.Type, PERSON) {
					s := globals.SENSITIVE
					f.SetContentCategory(&file.KazoupCategorization{
						ContentCategory: &s,
					})
					break
				}
			}

			n := time.Now()
			if tm == nil {
				f.SetOptsTimestamps(&file.OptsKazoupFile{
					TextAnalyzedTimestamp: &n,
				})
			} else {
				tm.TextAnalyzedTimestamp = &n
				f.SetOptsTimestamps(tm)
			}
		}

		b, err := json.Marshal(f)
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
	}

	return nil
}
