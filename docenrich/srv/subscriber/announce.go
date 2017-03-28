package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	doc_proto "github.com/kazoup/platform/docenrich/srv/proto/docenrich"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type AnnounceDocEnrich struct {
	Client client.Client
	Broker broker.Broker
}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceDocEnrich) OnCrawlerFinished(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled document files
	if globals.DiscoverTopic == msg.Handler {
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Call DocEnrich to kick off document enrichment over datasource files
		areq := a.Client.NewRequest(
			globals.DOCENRICH_SERVICE_NAME,
			"DocEnrich.Create",
			&doc_proto.CreateRequest{
				Type:  globals.TypeDatasource,
				Index: e.Index,
				Id:    e.Id,
			},
		)
		arsp := &doc_proto.CreateResponse{}

		if err := a.Client.Call(ctx, areq, arsp); err != nil {
			log.Println("ERROR Calling DocEnrich.Create", err)
		}
	}

	return nil
}
