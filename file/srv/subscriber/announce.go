package subscriber

import (
	"encoding/json"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/globals"
	announce_msg "github.com/kazoup/platform/lib/protomsg/announce"
	deletefile_msg "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type AnnounceFile struct {
	Client client.Client
	Broker broker.Broker
}

// OnFileDeleted
func (a *AnnounceFile) OnFileDeleted(ctx context.Context, msg *announce_msg.AnnounceMessage) error {
	// After file has been deleted, remove its thumbnail from our GCS account
	if globals.HANDLER_FILE_DELETE == msg.Handler {
		var r *file_proto.DeleteRequest
		if err := json.Unmarshal([]byte(msg.Data), &r); err != nil {
			return err
		}

		// Trigger deletion for associated resources (thumbnail in our GCS account)
		if err := a.Client.Publish(ctx, a.Client.NewPublication(globals.DeleteFileInBucketTopic, &deletefile_msg.DeleteFileInBucketMsg{
			FileId: r.FileId,
			Index:  r.Index,
		})); err != nil {
			return err
		}
	}

	return nil
}
