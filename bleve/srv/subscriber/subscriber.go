package subscriber

import (
	"github.com/blevesearch/bleve"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"golang.org/x/net/context"
)

type FileSubscriber struct {
	Index bleve.Index
}

func (fs *FileSubscriber) Handle(ctx context.Context, msg *crawler.FileMessage) error {
	//log.Printf("Got message: %s", msg.Say)

	//err := elastic.Bulk.Index("files", "file", "", "", "", nil, msg.Say)
	return nil
}
