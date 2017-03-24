package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	imgenrich_proto "github.com/kazoup/platform/imgenrich/srv/proto/imgenrich"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type AnnounceImgEnrich struct {
	Client client.Client
	Broker broker.Broker
}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceImgEnrich) OnCrawlerFinished(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled document files
	if globals.DiscoverTopic == msg.Handler {
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Call DocEnrich to kick off document enrichment over datasource files
		areq := a.Client.NewRequest(
			globals.IMGENRICH_SERVICE_NAME,
			"ImgEnrich.Create",
			&imgenrich_proto.CreateRequest{
				Type:  globals.TypeDatasource,
				Index: e.Index,
				Id:    e.Id,
			},
		)
		arsp := &imgenrich_proto.CreateResponse{}

		if err := a.Client.Call(ctx, areq, arsp); err != nil {
			log.Println("ERROR Calling ImgEnrich.Create for Datasource", err)
		}
	}

	return nil
}
