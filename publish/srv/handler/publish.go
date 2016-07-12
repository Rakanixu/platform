package handler

import (
	publish "github.com/kazoup/platform/publish/srv/proto/publish"
	"github.com/kazoup/platform/publish/srv/publisher"
	"golang.org/x/net/context"
)

// Publish ...
type Publish struct{}

// Send data within topic to be publish
func (p *Publish) Send(ctx context.Context, req *publish.SendRequest, res *publish.SendResponse) error {
	var header = make(map[string]string)

	header["topic"] = req.Topic
	header["id"] = req.Id
	header["publisher"] = req.Publisher

	err := publisher.PublishData(header, []byte(req.Data))

	if err != nil {
		// TODO: com.kazoup.srv.publish
		return err
	}

	return nil
}
