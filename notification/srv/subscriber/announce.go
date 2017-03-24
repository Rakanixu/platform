package subscriber

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	enrich_msg "github.com/kazoup/platform/lib/protomsg/enrich"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type AnnounceNotification struct {
	Client client.Client
	Broker broker.Broker
}

// OnDocEnrich
func (a *AnnounceNotification) OnDocEnrich(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Notify that document enrichment happened
	if globals.DocEnrichTopic == msg.Handler {
		var m *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		// Publish if requested
		if m.Notify {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
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
func (a *AnnounceNotification) OnImgEnrich(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Notify that image enrichment happened
	if globals.ImgEnrichTopic == msg.Handler {
		var m *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		// Publish if requested
		if m.Notify {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
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
func (a *AnnounceNotification) OnAudioEnrich(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Notify that audio enrichment happened
	if globals.AudioEnrichTopic == msg.Handler {
		var m *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		// Publish if requested
		if m.Notify {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
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
func (a *AnnounceNotification) OnSentimentExtraction(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Notify that sentiment extraction happened
	if globals.SentimentEnrichTopic == msg.Handler {
		var m *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		// Publish if requested
		if m.Notify {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
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
func (a *AnnounceNotification) OnEntitiesExtraction(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// Notify that entities extraction happened
	if globals.ExtractEntitiesTopic == msg.Handler {
		var m *enrich_msg.EnrichMessage
		if err := json.Unmarshal([]byte(msg.Data), &m); err != nil {
			return err
		}

		// Publish if requested
		if m.Notify {
			if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
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
