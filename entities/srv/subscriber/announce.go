package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/entities/srv/proto/entities"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnAudioEnrich reacts to AudioEnrichTopic
func (a *AnnounceHandler) OnAudioEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a an audio file has been enriched, we want to extract entities from content
	if globals.AudioEnrichTopic == msg.Handler {
		var e *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Extract entities
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ExtractEntitiesTopic, e)); err != nil {
			return err
		}
	}

	return nil
}

// OnDocEnrich reacts to DocEnrichTopic
func (a *AnnounceHandler) OnDocEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a an document file has been enriched, we want to extract entities from content
	if globals.DocEnrichTopic == msg.Handler {
		var e *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Extract entities
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ExtractEntitiesTopic, e)); err != nil {
			return err
		}
	}

	return nil
}

// OnExtractFile reacts to ExtractFile handler
func (a *AnnounceHandler) OnExtractFile(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_ENTITIES_EXTRACT_FILE == msg.Handler {
		var e *proto_entities.ExtractFileRequest
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		uID, err := utils.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ExtractEntitiesTopic, &enrich.EnrichMessage{
			Index:  e.Index,
			Id:     e.Id,
			UserId: uID,
			Notify: true,
		})); err != nil {
			return err
		}
	}

	return nil
}
