package subscriber

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	text "github.com/kazoup/platform/lib/normalization/text"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
	"time"
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
	ctx  context.Context
	msg  *enrich_proto.EnrichMessage
	done chan bool
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (t *taskHandler) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	c := enrichMsgChan{
		ctx:  ctx,
		msg:  enrichmsg,
		done: make(chan bool),
	}
	// Queue internally
	t.enrichMsgChan <- c

	<-c.done

	return nil
}

func (t *taskHandler) queueListener(wID int) {
	for m := range t.enrichMsgChan {
		if err := processEnrichMsg(m); err != nil {
			log.Println("Error Processing sentiment analyzer on worker", wID, err)
		}
	}
}

func startWorkers(t *taskHandler) {
	for i := 0; i < t.workers; i++ {
		go t.queueListener(i)
	}
}

func processEnrichMsg(m enrichMsgChan) error {
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

	sentimentTextAnalyzer := false
	tm := f.GetOptsTimestamps()

	if tm == nil || tm.SentimentAnalyzedTimestamp == nil {
		sentimentTextAnalyzer = true
	} else {
		sentimentTextAnalyzer = tm.SentimentAnalyzedTimestamp.Before(f.GetModifiedTime())
	}

	if sentimentTextAnalyzer {
		// Apply rossete sentiment
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

			s, err := rossetelib.Sentiment(t)
			if err != nil {
				return err
			}
			f.SetSentiment(s)

			n := time.Now()
			if tm == nil {
				f.SetOptsTimestamps(&file.OptsKazoupFile{
					SentimentAnalyzedTimestamp: &n,
				})
			} else {
				tm.SentimentAnalyzedTimestamp = &n
				f.SetOptsTimestamps(tm)
			}

			b, err := json.Marshal(f)
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
				log.Println("ERROR UPDATING FILE", err)
				return err
			}
		}

		m.msg.FileName = f.GetName()
	}

	m.done <- true

	return nil
}
