package handler

import (
	"log"

	"golang.org/x/net/context"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"

	elasticsearch "github.com/kazoup/platform/elastic/srv/proto/elastic"
)

func Subscriber(p broker.Publication) error {
	log.Printf("Got message: %v with id: %v", string(p.Message().Body), p.Message().Header["id"])
	ctx := context.TODO()
	msg := p.Message()
	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Create",
		&elasticsearch.CreateRequest{
			Index: "files", // Hardcoded index for flags
			Type:  "file",  // Hardcoded type ...
			Data:  string(msg.Body),
		},
	)
	srvRsp := &elasticsearch.CreateResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {

		log.Printf("ES response: %v %v ", srvRsp, err.Error())
		return errors.InternalServerError("go.micro.srv.indexer.Subscriber", err.Error())
	}
	return nil
}
