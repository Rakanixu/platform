package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	enrich_msg "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type AnnounceTextAnalyzer struct {
	Client client.Client
	Broker broker.Broker
}

// OnAudioEnrich reacts to AudioEnrichTopic
func (a *AnnounceTextAnalyzer) OnAudioEnrich(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a an audio file has been enriched, we want to extract entities from content
	if globals.AudioEnrichTopic == msg.Handler {
		var e *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Analyze file (Extract entities)
		if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.ExtractEntitiesTopic, e)); err != nil {
			return err
		}
	}

	return nil
}

// OnDocEnrich reacts to DocEnrichTopic
func (a *AnnounceTextAnalyzer) OnDocEnrich(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After a an document file has been enriched, we want to extract entities from content
	if globals.DocEnrichTopic == msg.Handler {
		var e *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Publish the same message to SentimentEnrichTopic
		if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.ExtractEntitiesTopic, e)); err != nil {
			return err
		}
	}

	return nil
}
