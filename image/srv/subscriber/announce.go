package subscriber

import (
	"encoding/json"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto_db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/image/srv/proto/image"
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
	// After a crawler has finished, we want enrich crawled document files
	if globals.DiscoverTopic == msg.Handler {
		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		if err := publishImgFilesNotProcessed(ctx, e); err != nil {
			return err
		}
	}

	return nil
}

// OnEnrichFile reacts to EnrichFile handler
func (a *AnnounceHandler) OnEnrichFile(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_IMAGE_ENRICH_FILE == msg.Handler {
		var r *proto_image.EnrichFileRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		uID, err := globals.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ImgEnrichTopic, &enrich.EnrichMessage{
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

// OnEnrichDatasource reacts to EnrichDatasource handler
func (a *AnnounceHandler) OnEnrichDatasource(ctx context.Context, msg *announce.AnnounceMessage) error {
	if globals.HANDLER_IMAGE_ENRICH_DATASOURCE == msg.Handler {
		var r *proto_image.EnrichDatasourceRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		rsp, err := db_helper.ReadFromDB(srv.Client(), ctx, &proto_db.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    r.Id,
		})
		if err != nil {
			return err
		}

		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(rsp.Result), &e); err != nil {
			return err
		}

		if err := publishImgFilesNotProcessed(ctx, e); err != nil {
			return err
		}
	}

	return nil
}
