package subscriber

import (
	"encoding/json"
	"fmt"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kennygrant/sanitize"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"regexp"
	"time"
)

type EnrichMsgChan struct {
	ctx context.Context
	msg *enrich_proto.EnrichMessage
}

type SentimentAnalyzer struct {
	Client        client.Client
	EnrichMsgChan chan EnrichMsgChan
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (sa *SentimentAnalyzer) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	// Queue internally
	sa.EnrichMsgChan <- EnrichMsgChan{
		ctx: ctx,
		msg: enrichmsg,
	}

	return nil
}

func (sa *SentimentAnalyzer) queueListener(wID int) {
	for m := range sa.EnrichMsgChan {
		if err := processEnrichMsg(sa.Client, m); err != nil {
			log.Println("Error Processing sentiment analyzer on worker", wID, err)
		}
	}
}

func StartWorkers(sa *SentimentAnalyzer) {
	for i := 0; i < sa.Workers; i++ {
		go sa.queueListener(i)
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

	sentimentTextAnalyzer := false
	tm := f.GetOptsTimestamps()

	if tm == nil {
		sentimentTextAnalyzer = true
	} else {
		sentimentTextAnalyzer = tm.SentimentAnalyzedTimestamp.Before(f.GetModifiedTime())
	}

	if sentimentTextAnalyzer {
		// Apply rossete sentiment
		if len(f.GetContent()) > 0 {
			nl, err := regexp.Compile("\n")
			if err != nil {
				return err
			}
			q, err := regexp.Compile("\"")
			if err != nil {
				return err
			}

			s, err := rossetelib.Sentiment(q.ReplaceAllString(nl.ReplaceAllString(sanitize.HTML(f.GetContent()), " "), ""))
			if err != nil {
				return err
			}
			f.SetSentiment(s)

			if tm == nil {
				f.SetOptsTimestamps(&file.OptsKazoupFile{
					SentimentAnalyzedTimestamp: time.Now(),
				})
			} else {
				tm.SentimentAnalyzedTimestamp = time.Now()
				f.SetOptsTimestamps(tm)
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
				log.Println("ERROR UPDATING FILE", err)
				return err
			}
		}

		// Publish notification topic if requested
		if m.msg.Notify {
			if err := c.Publish(m.ctx, c.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.msg.UserId,
				Info:   fmt.Sprintf("Sentiment extraction for %s finished.", f.GetName()),
			})); err != nil {
				log.Print("Publishing NotificationTopic (SentimentAnalyzer) error %s", err)
			}
		}
	}

	return nil
}
