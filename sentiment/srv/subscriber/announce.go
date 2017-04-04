package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/sentiment/srv/proto/sentiment"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type SentimentHandler struct{}

// OnAudioEnrich reacts to AudioEnrichTopic
func (s *SentimentHandler) OnAudioEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a an audio file has been enriched, we want to extract entities from content
	if globals.AudioEnrichTopic == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var e *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Analyze file (Extract entities)
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.SentimentEnrichTopic, e)); err != nil {
			return err
		}
	}

	return nil
}

// OnDocEnrich reacts to DocEnrichTopic
func (s *SentimentHandler) OnDocEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a an document file has been enriched, we want to extract entities from content
	if globals.DocEnrichTopic == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var e *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		// Publish the same message to SentimentEnrichTopic
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.SentimentEnrichTopic, e)); err != nil {
			return err
		}
	}

	return nil
}

// OnAnalyzeFile reacts to AnalyzeFile handler
func (s *SentimentHandler) OnAnalyzeFile(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_SENTIMENT_ENRICH_FILE == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var r *proto_sentiment.AnalyzeFileRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		uID, err := globals.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.SentimentEnrichTopic, &enrich.EnrichMessage{
			Index:  r.Index,
			Id:     r.Id,
			UserId: uID,
			Notify: true,
		})); err != nil {
			return err
		}
	}

	return nil
}
