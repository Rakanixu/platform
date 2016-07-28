package subscriber

import (
	"log"

	"github.com/kazoup/platform/elastic/srv/elastic"
	example "github.com/micro/micro/examples/template/srv/proto/example"
	"golang.org/x/net/context"
)

func FileSubscriber(ctx context.Context, msg *example.Message) error {
	//log.Printf("Got message: %s", msg.Say)

	err := elastic.Bulk.Index("files", "file", "", "", "", nil, msg.Say)
	if err != nil {
		log.Print("Bulk Indexer error %s", err.Error())
		return err
	}

	return nil
}
