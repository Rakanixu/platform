package subscriber

import (
	model "github.com/kazoup/platform/db/srv/engine/elastic/model"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"golang.org/x/net/context"
	elib "gopkg.in/olivere/elastic.v5"
	"log"
)

func Subscribe(e *model.Elastic) error {

	// Files
	go func() {
		for {
			select {
			case v := <-e.FilesChannel:
				// File message can be notified, when a file is create, deleted or shared within kazoup
				if v.FileMessage.Notify {
					// We do not use bulk, as is just one element
					_, err := e.Client.Index().Index(v.FileMessage.Index).Type(globals.FileType).Id(v.FileMessage.Id).BodyString(v.FileMessage.Data).Do(context.Background())
					if err != nil {
						log.Print("Indexer error %s", err)
					}

					n := &notification_proto.NotificationMessage{
						Method: globals.NOTIFY_REFRESH_SEARCH,
						UserId: v.FileMessage.UserId,
					}

					// Publish scan topic, crawlers should pick up message
					if err := v.Client.Publish(globals.NewSystemContext(), v.Client.NewPublication(globals.NotificationTopic, n)); err != nil {
						log.Print("Publishing (notify file) error %s", err)
					}
				} else {
					// Use bulk processor as we will index groups of documents
					r := elib.NewBulkIndexRequest().Index(v.FileMessage.Index).Type(globals.FileType).Id(v.FileMessage.Id).Doc(v.FileMessage.Data)
					e.BulkProcessor.Add(r)
				}
			}

		}
	}()

	// Slack users
	go func() {
		for {
			select {
			case v := <-e.SlackUsersChannel:
				// Use bulk processor as we will index groups of documents
				r := elib.NewBulkIndexRequest().Index(v.Index).Type(globals.UserType).Id(v.Id).Doc(v.Data)
				e.BulkProcessor.Add(r)
			}

		}
	}()

	// Slack channels
	go func() {
		for {
			select {
			case v := <-e.SlackChannelsChannel:
				// Use bulk processor as we will index groups of documents
				r := elib.NewBulkIndexRequest().Index(v.Index).Type(globals.ChannelType).Id(v.Id).Doc(v.Data)
				e.BulkProcessor.Add(r)
			}

		}
	}()

	return nil
}
