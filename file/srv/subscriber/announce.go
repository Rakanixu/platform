package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	deletefile "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// OnFileDeleted
func (a *AnnounceHandler) OnFileDeleted(ctx context.Context, msg *announce.AnnounceMessage) error {
	// After file has been deleted, remove its thumbnail from our GCS account
	if globals.HANDLER_FILE_DELETE == msg.Handler {
		var r *proto_file.DeleteRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Trigger deletion for associated resources (thumbnail in our GCS account)
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.DeleteFileInBucketTopic, &deletefile.DeleteFileInBucketMsg{
			FileId: r.FileId,
			Index:  r.Index,
		})); err != nil {
			return err
		}
	}

	return nil
}
