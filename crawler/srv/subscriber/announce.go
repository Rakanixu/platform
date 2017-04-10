package subscriber

import (
	"encoding/json"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	crawler "github.com/kazoup/platform/lib/protomsg/crawler"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnDiscover
func (a *AnnounceHandler) OnDiscover(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a scan finishes, we want to procced with post scan steeps
	if globals.DiscoverTopic == msg.Handler {
		var e *proto_datasource.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
			globals.DiscoveryFinishedTopic,
			&crawler.CrawlerFinishedMessage{
				DatasourceId: e.Id,
			},
		)); err != nil {
			return err
		}
	}

	return nil
}
