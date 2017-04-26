package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

const (
	AUDIO_TIMESTAMP_ES_MAP = "opts_kazoup_file.audio_timestamp"
)

type AnnounceHandler struct{}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceHandler) OnCrawlerFinished(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled audio files
	if globals.DiscoverTopic == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		rsp, err := custom.ScrollUnprocessedFiles(ctx, &proto_custom.ScrollUnprocessedFilesRequest{
			Index:    e.Index,
			Category: globals.CATEGORY_AUDIO,
			Field:    AUDIO_TIMESTAMP_ES_MAP,
		})
		if err != nil {
			return err
		}

		var r []*file.KazoupFile
		if err := json.Unmarshal([]byte(rsp.Result), &r); err != nil {
			return err
		}

		// Publish msg for all files not being process yet
		for _, v := range r {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.AudioEnrichTopic, &enrich.EnrichMessage{
				Index:  e.Index,
				Id:     v.ID,
				UserId: e.UserId,
			})); err != nil {
				log.Println("ERROR publishing AudioEnrichTopic", err)
			}
		}
	}

	return nil
}

func (a *AnnounceHandler) OnEnrichFile(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_AUDIO_ENRICH_FILE == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		uID, err := globals.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		var req *proto_audio.EnrichFileRequest
		if err := json.Unmarshal([]byte(msg.Data), &req); err != nil {
			return err
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.AudioEnrichTopic, &enrich.EnrichMessage{
			Index:  req.Index,
			Id:     req.Id,
			UserId: uID,
			Notify: true,
		})); err != nil {
			return err
		}
	}

	return nil
}

func (a *AnnounceHandler) OnEnrichDatasource(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_AUDIO_ENRICH_DATASOURCE == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var req *proto_audio.EnrichDatasourceRequest
		if err := json.Unmarshal([]byte(msg.Data), &req); err != nil {
			return err
		}

		rsp, err := operations.Read(ctx, &proto_operations.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    req.Id,
		})
		if err != nil {
			return err
		}

		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(rsp.Result), &e); err != nil {
			return err
		}

		srsp, err := custom.ScrollUnprocessedFiles(ctx, &proto_custom.ScrollUnprocessedFilesRequest{
			Index:    e.Index,
			Category: globals.CATEGORY_AUDIO,
			Field:    AUDIO_TIMESTAMP_ES_MAP,
		})
		if err != nil {
			return err
		}

		var r []*file.KazoupFile
		if err := json.Unmarshal([]byte(srsp.Result), &r); err != nil {
			return err
		}

		// Publish msg for all files not being process yet
		for _, v := range r {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.AudioEnrichTopic, &enrich.EnrichMessage{
				Index:  e.Index,
				Id:     v.ID,
				UserId: e.UserId,
			})); err != nil {
				log.Println("ERROR publishing AudioEnrichTopic", err)
			}
		}
	}

	return nil
}
