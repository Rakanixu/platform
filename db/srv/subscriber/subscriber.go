package subscriber

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"golang.org/x/net/context"
)

type FileSubscriber struct {
	MsgCh chan *crawler.FileMessage
}

func (fs *FileSubscriber) Handle(ctx context.Context, msg *crawler.FileMessage) error {
	fs.MsgCh <- msg
	return nil
}
