package subscriber

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	enrich_proto "github.com/kazoup/platform/enrich/srv/proto/enrich"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type Enrich struct {
	Client client.Client
}

// Enrich subscriber, receive EnrichMessage to get the file and process it
func (e *Enrich) Enrich(ctx context.Context, enrichmsg enrich_proto.EnrichMessage) error {
	rsp, err := db_helper.ReadFromDB(e.Client, context.Background(), &db_proto.ReadRequest{
		Index: enrichmsg.Index,
		Type:  globals.FileType,
		Id:    enrichmsg.Id,
	})
	if err != nil {
		return err
	}

	f, err := file.NewFileFromString(rsp.Result)
	if err != nil {
		return err
	}

	log.Println("FILE YEIII", f)

	return nil
}
