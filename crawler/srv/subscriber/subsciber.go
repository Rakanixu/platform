package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	crawler_proto "github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/fs"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

// Scans subscriber, receive endpoint to crawl it
func Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	fs, err := fs.NewFsFromEndpoint(endpoint)
	if err != nil {
		return err
	}

	// Publish crawler started, or is just going to start..
	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerStartedTopic, &crawler_proto.CrawlerStartedMessage{
		UserId:       endpoint.UserId,
		DatasourceId: endpoint.Id,
	})); err != nil {
		return err
	}

	// Receive files founded by FileSystem
	c, r, err := fs.List()
	if err != nil {
		return err
	}

	for {
		select {
		// Channel receives signal cralwer has finished
		case <-r:
			time.Sleep(time.Second * 5)
			// Clear index (files that no longer exists, rename, etc..)
			if err := globals.ClearIndex(endpoint); err != nil {
				log.Println("ERROR clearing index after scan", err)
			}

			msg := &crawler.CrawlerFinishedMessage{
				DatasourceId: endpoint.Id,
			}
			// Publish crawling process has finished
			if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
				return err
			}
			close(c)
			close(r)

			return nil
		// Channel receives File to be indexed by Elastic Search
		case f := <-c:
			if err := file.IndexAsync(f, globals.FilesTopic, f.GetIndex(), false); err != nil {
				log.Println("Error indexing async file", err)
			}
		}
	}

	return nil
}
