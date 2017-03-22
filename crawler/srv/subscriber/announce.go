package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	cawler_msg "github.com/kazoup/platform/lib/protomsg/crawler"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type AnnounceCrawler struct {
	Client client.Client
	Broker broker.Broker
}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceCrawler) OnCrawlerFinished(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a crawler has finished, we want to notify user
	if globals.ScanTopic == msg.Handler {
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		if err := a.Client.Publish(ctx, a.Client.NewPublication(
			globals.CrawlerFinishedTopic,
			&cawler_msg.CrawlerFinishedMessage{
				DatasourceId: e.Id,
			},
		)); err != nil {
			return err
		}
	}

	return nil
}
