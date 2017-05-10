package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	crawler "github.com/kazoup/platform/lib/protomsg/crawler"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnDocEnrich
func (a *AnnounceHandler) OnDocEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Notify that document enrichment happened
	if globals.DocEnrichTopic == msg.Handler {
		var m *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Publish if requested
		if m.Notify {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.UserId,
				Info:   fmt.Sprintf("Document content extraction for %s finished.", m.FileName),
			})); err != nil {
				return err
			}
		}
	}

	return nil
}

// OnImgEnrich
func (a *AnnounceHandler) OnImgEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Notify that image enrichment happened
	if globals.ImgEnrichTopic == msg.Handler {
		var m *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Publish if requested
		if m.Notify {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.UserId,
				Info:   fmt.Sprintf("Image content extraction for %s finished.", m.FileName),
			})); err != nil {
				return err
			}
		}
	}

	return nil
}

// OnAudioEnrich
func (a *AnnounceHandler) OnAudioEnrich(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Notify that audio enrichment happened
	if globals.AudioEnrichTopic == msg.Handler {
		var m *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Publish if requested
		if m.Notify {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.UserId,
				Info:   fmt.Sprintf("Speach to text for %s finished.", m.FileName),
			})); err != nil {
				return err
			}
		}
	}

	return nil
}

// OnSentimentExtraction
func (a *AnnounceHandler) OnSentimentExtraction(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Notify that sentiment extraction happened
	if globals.SentimentEnrichTopic == msg.Handler {
		var m *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Publish if requested
		if m.Notify {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.UserId,
				Info:   fmt.Sprintf("Sentiment extraction for %s finished.", m.FileName),
			})); err != nil {
				return err
			}
		}
	}

	return nil
}

// OnSentimentExtraction
func (a *AnnounceHandler) OnEntitiesExtraction(ctx context.Context, msg *announce.AnnounceMessage) error {
	// Notify that entities extraction happened
	if globals.ExtractEntitiesTopic == msg.Handler {
		var m *enrich.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Publish if requested
		if m.Notify {
			if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
				Method: globals.NOTIFY_REFRESH_SEARCH,
				UserId: m.UserId,
				Info:   fmt.Sprintf("Entity extraction for %s finished.", m.FileName),
			})); err != nil {
				return err
			}
		}

	}

	return nil
}

//OnCrawlerFinished
func (a *AnnounceHandler) OnCrawlerFinished(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a crawler finishes, we want to notify user
	if globals.DiscoveryFinishedTopic == msg.Handler {
		var m *crawler.CrawlerFinishedMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		rsp, err := operations.Read(ctx, &proto_operations.ReadRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    m.DatasourceId,
		})
		if err != nil {
			return err
		}

		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(rsp.Result), &e); err != nil {
			return err
		}

		// Publish notification
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
			Info:   "Scan finished on " + e.Url + " datasource.",
			Method: globals.NOTIFY_REFRESH_DATASOURCES,
			UserId: e.UserId,
			Data:   rsp.Result,
		})); err != nil {
			return err
		}
	}

	return nil
}

// OnFileDeleted
func (a *AnnounceHandler) OnFileDeleted(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After file has been deleted, remove its thumbnail from our GCS account
	if globals.HANDLER_FILE_DELETE == msg.Handler {
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Get userId
		uId, err := globals.ParseUserIdFromContext(ctx)
		if err != nil {
			return err
		}

		// Publish notification
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.NotificationTopic, &proto_notification.NotificationMessage{
			Method: globals.NOTIFY_REFRESH_SEARCH,
			UserId: uId,
		})); err != nil {
			return err
		}
	}

	return nil
}
