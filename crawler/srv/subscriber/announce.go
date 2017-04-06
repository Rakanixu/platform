package subscriber

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	cawler_msg "github.com/kazoup/platform/lib/protomsg/crawler"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnScanFinished
func (a *AnnounceHandler) OnScanFinished(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After a scan finishes, we want to procced with post scan steeps
	if globals.DiscoverTopic == msg.Handler {
		var e *datasource_proto.Endpoint
		if err := json.Unmarshal([]byte(msg.Data), &e); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(
			globals.CrawlerFinishedTopic,
			&cawler_msg.CrawlerFinishedMessage{
				DatasourceId: e.Id,
			},
		)); err != nil {
			return err
		}
	}

	return nil
}
