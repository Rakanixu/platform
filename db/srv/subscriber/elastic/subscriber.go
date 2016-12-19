package subscriber

import (
	model "github.com/kazoup/platform/db/srv/engine/elastic/model"
	"github.com/kazoup/platform/lib/globals"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
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
					if _, err := e.Conn.Index(v.FileMessage.Index, "file", v.FileMessage.Id, nil, v.FileMessage.Data); err != nil {
						log.Print("Bulk Indexer error %s", err)
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
					// Use bulk as we will index groups of documents
					if err := e.Bulk.Index(v.FileMessage.Index, "file", v.FileMessage.Id, "", "", nil, v.FileMessage.Data); err != nil {
						log.Print("Bulk Indexer error %s", err)
					}
				}
			}

		}
	}()

	// Slack users
	go func() {
		for {
			select {
			case v := <-e.SlackUsersChannel:
				if err := e.Bulk.Index(v.Index, "user", v.Id, "", "", nil, v.Data); err != nil {
					log.Print("Bulk Indexer error %s", err)
				}
			}

		}
	}()

	// Slack channels
	go func() {
		for {
			select {
			case v := <-e.SlackChannelsChannel:
				if err := e.Bulk.Index(v.Index, "channel", v.Id, "", "", nil, v.Data); err != nil {
					log.Print("Bulk Indexer error %s", err)
				}
			}

		}
	}()

	return nil
}
