package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/image/srv/proto/image"
	"github.com/kazoup/platform/lib/db/custom"
	"github.com/kazoup/platform/lib/db/custom/proto/custom"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

const (
	IMAGE_TIMESTAMP_ES_MAP = "opts_kazoup_file.tags_timestamp"
)

type AnnounceHandler struct{}

// OnCrawlerFinished reacts to CrawlerFinished
func (a *AnnounceHandler) OnCrawlerFinished(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a crawler has finished, we want enrich crawled document files
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
			Category: globals.CATEGORY_PICTURE,
			Field:    IMAGE_TIMESTAMP_ES_MAP,
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
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ImgEnrichTopic, &enrich.EnrichMessage{
				Index:  e.Index,
				Id:     v.ID,
				UserId: e.UserId,
			})); err != nil {
				log.Println("ERROR publishing ImgEnrichTopic", err)
			}
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

		uID, err := utils.ParseUserIdFromContext(ctx)
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
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		var r *proto_image.EnrichDatasourceRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		rrsp, err := operations.Read(ctx, &proto_operations.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    r.Id,
		})
		if err != nil {
			return err
		}

		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(rrsp.Result), &e); err != nil {
			return err
		}

		srsp, err := custom.ScrollUnprocessedFiles(ctx, &proto_custom.ScrollUnprocessedFilesRequest{
			Index:    e.Index,
			Category: globals.CATEGORY_PICTURE,
			Field:    IMAGE_TIMESTAMP_ES_MAP,
		})
		if err != nil {
			return err
		}

		var kf []*file.KazoupFile
		if err := json.Unmarshal([]byte(srsp.Result), &kf); err != nil {
			return err
		}

		// Publish msg for all files not being process yet
		for _, v := range kf {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.ImgEnrichTopic, &enrich.EnrichMessage{
				Index:  e.Index,
				Id:     v.ID,
				UserId: e.UserId,
			})); err != nil {
				log.Println("ERROR publishing ImgEnrichTopic", err)
			}
		}
	}

	return nil
}
