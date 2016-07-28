package handler

import (
	"golang.org/x/net/context"

	elasticsearch "github.com/kazoup/platform/elastic/srv/proto/elastic"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	example "github.com/micro/micro/examples/template/srv/proto/example"
)

func FileSubscriber(ctx context.Context, msg *example.Message) error {
	//log.Printf("Got message: %s", msg.Say)

	ctx = context.TODO()
	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.BulkCreate",
		&elasticsearch.BulkCreateRequest{
			Index: "files", // Hardcoded index for flags
			Type:  "file",  // Hardcoded type ...
			Data:  string(msg.Say),
		},
	)
	srvRsp := &elasticsearch.BulkCreateResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {

		return errors.InternalServerError("go.micro.srv.indexer.Subscriber", err.Error())
	}
	return nil
}
