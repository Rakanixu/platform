package subscriber

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"github.com/kennygrant/sanitize"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"regexp"
	"time"
)

type TextAnalyzer struct {
	Client        client.Client
	EnrichMsgChan chan *enrich_proto.EnrichMessage
	Workers       int
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (ta *TextAnalyzer) Enrich(ctx context.Context, enrichmsg *enrich_proto.EnrichMessage) error {
	// Queue internally
	ta.EnrichMsgChan <- enrichmsg

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
			nl, err := regexp.Compile("\n")
			if err != nil {
				return err
			}
			q, err := regexp.Compile("\"")
			if err != nil {
				return err
			}

			e, err := rossetelib.Entities(q.ReplaceAllString(nl.ReplaceAllString(sanitize.HTML(f.GetContent()), " "), ""))
			if err != nil {
				return err
			}
			f.SetEntities(e)

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

		_, err = db_helper.UpdateFromDB(c, globals.NewSystemContext(), &db_proto.UpdateRequest{
			Index: m.Index,
			Type:  globals.FileType,
			Id:    m.Id,
			Data:  string(b),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
