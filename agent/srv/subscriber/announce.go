package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	save "github.com/kazoup/platform/lib/protomsg/saveremotefiles"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type AnnounceHandler struct{}

// Checks if announcement sent to the general topic is for the agent service, if yes,
// resends the message to another subscriber that logs it
func (a *AnnounceHandler) OnSave(ctx context.Context, msg *announce.AnnounceMessage) error {
	// If the handler is agent handler
	if globals.HANDLER_AGENT_SAVE == msg.Handler {
		// Get server from context
		srv, ok := micro.FromContext(ctx)
		if !ok {
			return errors.ErrInvalidCtx
		}

		// Unmarshal Kazoup file
		var f file.KazoupFile
		err := json.Unmarshal([]byte(msg.Data), &f)
		if err != nil {
			return err
		}

		// Publish new message to SaveRemoteFiles topic
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.SaveRemoteFilesTopic, &save.SaveRemoteFilesMessage{
			UserId: f.UserId,
			Index:  f.Index,
			Data:   msg.Data,
		})); err != nil {
			return err
		}
	}

	return nil
}
