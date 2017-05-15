package elastic

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/db/bulk"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/lib/slack"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"log"
	"os"
	"time"
)

type elastic struct {
	Client               *elib.Client
	BulkProcessor        *elib.BulkProcessor
	BulkFilesProcessor   *elib.BulkProcessor
	FilesChannel         chan *filesChannel
	SlackUsersChannel    chan *crawler.SlackUserMessage
	SlackChannelsChannel chan *crawler.SlackChannelMessage
}

type filesChannel struct {
	FileMessage *crawler.FileMessage
	Ctx         context.Context
}

type bulkableKazoupRequest struct {
	context.Context
	elib.BulkableRequest
}

type processData func(e *elastic)

func init() {
	bulk.Register(&elastic{
		FilesChannel:         make(chan *filesChannel),
		SlackUsersChannel:    make(chan *crawler.SlackUserMessage),
		SlackChannelsChannel: make(chan *crawler.SlackChannelMessage),
	})
}

// Init Elastic db (engine)
// Common init for DB, Config and Subscriber interfaces
func (e *elastic) Init(srv micro.Service) error {
	var err error

	// Set ES details from env variables
	url := os.Getenv("ELASTICSEARCH_URL")
	if url == "" {
		url = "http://elasticsearch:9200"
	}
	username := os.Getenv("ES_USERNAME")
	password := os.Getenv("ES_PASSWORD")

	// Client
	e.Client, err = elib.NewSimpleClient(
		elib.SetURL(url),
		elib.SetBasicAuth(username, password),
		elib.SetMaxRetries(3),
	)
	if err != nil {
		return err
	}

	rs, err := utils.NewUUID()
	if err != nil {
		return err
	}

	// Bulk Processor, used for users and channels
	e.BulkProcessor, err = e.Client.BulkProcessor().
		Name(fmt.Sprintf("bulkProcessor-%s", rs)).
		Workers(3).
		BulkActions(100).                // commit if # requests >= 100
		BulkSize(2 << 20).               // commit if size of requests >= 2 MB, probably to big, btw other constrains will be hit before
		FlushInterval(10 * time.Second). // commit every 10s
		Do(context.Background())
	if err != nil {
		return err
	}

	// Bulk Files Processor, used for index After function to
	e.BulkFilesProcessor, err = e.Client.BulkProcessor().
		After(func(executionId int64, requests []elib.BulkableRequest, response *elib.BulkResponse, err error) {
			for _, req := range requests {
				type updateBody struct {
					Doc *file.KazoupFile `json:"doc"`
				}

				var kf updateBody

				// elib.BulkableRequest stores two objects, headers and body
				src, err := req.Source()
				if err != nil {
					log.Printf("Error: %v", err)
					return
				}

				if len(src) == 2 {
					json.Unmarshal([]byte(src[1]), &kf)
				}

				// Assert type and use the proper context
				bkr, ok := req.(bulkableKazoupRequest)
				if !ok {
					log.Printf("Error bulkableKazoupRequest assertion: %v", bkr)
					return
				}

				if kf.Doc.Category == globals.CATEGORY_PICTURE &&
					(kf.Doc.OptsKazoupFile == nil || kf.Doc.OptsKazoupFile.ThumbnailTimestamp == nil) {
					if err := srv.Client().Publish(bkr.Context, srv.Client().NewPublication(globals.ThumbnailTopic, &enrich_proto.EnrichMessage{
						Index:  kf.Doc.Index,
						Id:     kf.Doc.ID,
						UserId: kf.Doc.UserId,
					})); err != nil {
						log.Printf("Publishing ThumbnailTopic error %s", err)
					}
				}
			}
		}).
		Name(fmt.Sprintf("bulkFilesProcessor-%s", rs)).
		Workers(3).
		BulkActions(500).               // commit if # requests >= 500
		BulkSize(10 << 20).             // commit if size of requests >= 10 MB, probably to big, btw other constrains will be hit before
		FlushInterval(6 * time.Second). // commit every 5s, notification message can be send and until 5s later is not really finished
		Do(context.Background())
	if err != nil {
		return err
	}

	// Files
	go processFiles(e)
	// Slack users
	go processSlackUsers(e)
	// Slack channels
	go processSlackChannels(e)

	return nil
}

// Files receive file messages to be indexed
func (e *elastic) Files(ctx context.Context, msg *crawler.FileMessage) error {
	e.FilesChannel <- &filesChannel{
		FileMessage: msg,
		Ctx:         ctx,
	}

	return nil
}

// SlackUsers receives slack user messages to be indexed
func (e *elastic) SlackUsers(ctx context.Context, msg *crawler.SlackUserMessage) error {
	e.SlackUsersChannel <- msg

	return nil
}

// SlackChannels receives slack channel messages to be indexed
func (e *elastic) SlackChannels(ctx context.Context, msg *crawler.SlackChannelMessage) error {
	e.SlackChannelsChannel <- msg

	return nil
}

func processFiles(e *elastic) {
	for {
		select {
		case v := <-e.FilesChannel:
			// File message can be notified, when a file is create, deleted or shared within kazoup
			if v.FileMessage.Notify {
				// We do not use bulk, as is just one element
				_, err := e.Client.Index().Index(v.FileMessage.Index).Type(globals.FileType).Id(v.FileMessage.Id).BodyString(v.FileMessage.Data).Do(context.Background())
				if err != nil {
					log.Printf("Indexer error %s", err)
				}

				srv, ok := micro.FromContext(v.Ctx)
				if !ok {
					log.Printf("%v", errors.ErrInvalidCtx)
				}

				n := &proto_notification.NotificationMessage{
					Method: globals.NOTIFY_REFRESH_SEARCH,
					UserId: v.FileMessage.UserId,
				}

				// Publish scan topic, crawlers should pick up message
				if err := srv.Client().Publish(v.Ctx, srv.Client().NewPublication(globals.NotificationTopic, n)); err != nil {
					log.Printf("Publishing (notify file) error %s", err)
				}
			} else {
				f, err := file.NewFileFromString(v.FileMessage.Data)
				if err != nil {
					log.Printf("Error creating file from string error %s", err)
				}

				// Use bulk processor as we will index groups of documents
				// We need to build the file to be able to update JSON like obj, and not string
				// We use helper BulkableKazoupRequest interface to do  not lose context on the after commit function
				r := bulkableKazoupRequest{
					v.Ctx,
					elib.NewBulkUpdateRequest().Index(v.FileMessage.Index).Type(globals.FileType).Id(v.FileMessage.Id).DocAsUpsert(true).Doc(f),
				}

				e.BulkFilesProcessor.Add(r)
			}
		}
	}
}

func processSlackUsers(e *elastic) {
	for {
		select {
		case v := <-e.SlackUsersChannel:
			var u slack.ESSlackUser
			if err := json.Unmarshal([]byte(v.Data), &u); err != nil {
				log.Printf("Error unmarshalling user %s", err)
			}

			// Use bulk processor as we will index groups of documents
			r := elib.NewBulkUpdateRequest().Index(v.Index).Type(globals.UserType).Id(utils.GetMD5Hash(u.UserID)).DocAsUpsert(true).Doc(u)
			e.BulkProcessor.Add(r)
		}

	}
}

func processSlackChannels(e *elastic) {
	for {
		select {
		case v := <-e.SlackChannelsChannel:
			var ch slack.ESSlackChannel
			if err := json.Unmarshal([]byte(v.Data), &ch); err != nil {
				log.Printf("Error unmarshalling channel %s", err)
			}

			// Use bulk processor as we will index groups of documents
			r := elib.NewBulkUpdateRequest().Index(v.Index).Type(globals.ChannelType).Id(utils.GetMD5Hash(ch.ChannelID)).DocAsUpsert(true).Doc(ch)
			e.BulkProcessor.Add(r)
		}

	}
}
