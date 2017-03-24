package subscriber

import (
	"encoding/json"
	audio_proto "github.com/kazoup/platform/audioenrich/srv/proto/audioenrich"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type AnnounceAudioEnrich struct {
	Client client.Client
	Broker broker.Broker
}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceAudioEnrich) OnCrawlerFinished(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled audio files
	if globals.DiscoverTopic == msg.Handler {
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Call AudioEnrich.Create to process datasource
		areq := a.Client.NewRequest(
			globals.AUDIOENRICH_SERVICE_NAME,
			"AudioEnrich.Create",
			&audio_proto.CreateRequest{
				Type:  globals.TypeDatasource,
				Index: e.Index,
				Id:    e.Id,
			},
		)
		arsp := &audio_proto.CreateResponse{}

		if err := a.Client.Call(ctx, areq, arsp); err != nil {
			log.Println("ERROR Calling AudioEnrich.Create for Datasource", err)
		}
	}

	return nil
}
