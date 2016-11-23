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

type Crawler struct {
	Client client.Client
}

// Scans subscriber, receive endpoint to crawl it
func (c *Crawler) Scans(ctx context.Context, endpoint *datasource.Endpoint) error {
	cfs, err := fs.NewFsFromEndpoint(endpoint)
	if err != nil {
		return err
	}

	// Publish crawler started, or is just going to start..
	if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.CrawlerStartedTopic, &crawler_proto.CrawlerStartedMessage{
		UserId:       endpoint.UserId,
		DatasourceId: endpoint.Id,
	})); err != nil {
		return err
	}

	// Receive files founded by FileSystem
	fc, r, err := cfs.List()
	if err != nil {
		return err
	}

	for {
		select {
		// Channel receives signal cralwer has finished
		case <-r:
			time.Sleep(time.Second * 8)

			// Clear index (files that no longer exists, rename, etc..)
			if err := globals.ClearIndex(c.Client, endpoint); err != nil {
				log.Println("ERROR clearing index after scan", err)
			}

			msg := &crawler.CrawlerFinishedMessage{
				DatasourceId: endpoint.Id,
			}
			// Publish crawling process has finished
			if err := c.Client.Publish(context.Background(), c.Client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
				return err
			}
			close(fc)
			close(r)

			return nil
		// Channel receives File to be indexed by Elastic Search
		case f := <-fc:
			if err := file.IndexAsync(c.Client, f, globals.FilesTopic, f.GetIndex(), false); err != nil {
				log.Println("Error indexing async file", err)
			}
		}
	}

	return nil
}
