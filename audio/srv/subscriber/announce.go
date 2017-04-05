package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto_db "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceHandler) OnCrawlerFinished(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled audio files
	if globals.DiscoverTopic == msg.Handler {
		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		if err := publishAudioFilesNotProcessed(ctx, e); err != nil {
			return err
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

		rsp, err := db_helper.ReadFromDB(srv.Client(), ctx, &proto_db.ReadRequest{
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

		if err := publishAudioFilesNotProcessed(ctx, e); err != nil {
			return err
		}
	}

	return nil
}
