package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/fs"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

func Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	fs, err := fs.NewFsFromEndpoint(endpoint)
	if err != nil {
		return err
	}

	c, r, err := fs.List()
	if err != nil {
		return err
	}

	// Publish notification
	msg := &notification_proto.NotificationMessage{
		Info: "Scan started on " + endpoint.Url + " datasource.",
	}

	if err := client.Publish(ctx, client.NewPublication(globals.NotificationTopic, msg)); err != nil {
		return err
	}

	for {
		select {
		case <-r:
			time.Sleep(time.Second * 5)

			if err := globals.ClearIndex(endpoint); err != nil {
				log.Println("ERROR clearing index after scan", err)
			}

			msg := &crawler.CrawlerFinishedMessage{
				DatasourceId: endpoint.Id,
			}

			if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
				return err
			}
			close(c)
			close(r)

			return nil
		case f := <-c:
			if err := file.IndexAsync(f, globals.FilesTopic, f.GetIndex()); err != nil {
				log.Println("Error indexing async file")
			}
		}
	}

	return nil
}
