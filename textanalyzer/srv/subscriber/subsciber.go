package subscriber

import (
	"encoding/json"
	"fmt"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	text "github.com/kazoup/platform/lib/normalization/text"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

const (
	//https://developer.rosette.com/features-and-functions#-entity-types
	IDENTIFIER = "IDENTIFIER"
	PERSON     = "PERSON"
)

type EnrichMsgChan struct {
	ctx context.Context
	msg *enrich_proto.EnrichMessage
}

type TextAnalyzer struct {
	Client        client.Client
	EnrichMsgChan chan EnrichMsgChan
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (ta *TextAnalyzer) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	// Queue internally
	ta.EnrichMsgChan <- EnrichMsgChan{
		ctx: ctx,
		msg: enrichmsg,
	}

	return nil
}

func (ta *TextAnalyzer) queueListener(wID int) {
	for m := range ta.EnrichMsgChan {
		if err := processEnrichMsg(ta.Client, m); err != nil {
			log.Println("Error Processing text analyzer on worker", wID, err)
		}
	}
}

func StartWorkers(ta *TextAnalyzer) {
	for i := 0; i < ta.Workers; i++ {
		go ta.queueListener(i)
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

	processTextAnalyzer := false
	tm := f.GetOptsTimestamps()

	if tm == nil {
		processTextAnalyzer = true
	} else {
		processTextAnalyzer = tm.TextAnalyzedTimestamp.Before(f.GetModifiedTime())
	}

	if processTextAnalyzer {
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

			e, err := rossetelib.Entities(t)
			if err != nil {
				return err
			}
			f.SetEntities(e)

			for _, ent := range e.Entities {
				if strings.Contains(ent.Type, IDENTIFIER) || strings.Contains(ent.Type, PERSON) {
					f.SetContentCategory(globals.SENSITIVE)
					break
				}
			}

			if tm == nil {
				f.SetOptsTimestamps(&file.OptsKazoupFile{
					TextAnalyzedTimestamp: time.Now(),
				})
			} else {
				tm.TextAnalyzedTimestamp = time.Now()
				f.SetOptsTimestamps(tm)
			}
		}

		b, err := json.Marshal(f)
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
				Info:   fmt.Sprintf("Entity extraction for %s finished.", f.GetName()),
			})); err != nil {
				log.Print("Publishing NotificationTopic (ImgEnrich) error %s", err)
			}
		}
	}

	return nil
}
