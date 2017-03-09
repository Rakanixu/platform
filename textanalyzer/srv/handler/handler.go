package handler

import (
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	proto "github.com/kazoup/platform/textanalyzer/srv/proto/textanalyzer"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type TextAnalyzer struct {
	Client client.Client
}

// Create handler publish a SINGLE FILE ExtractEntitiesTopic msg (in difference of imgenrich, docenrich, this endpoint DOES NOT SUPPORT DATASOURCE)
func (ta *TextAnalyzer) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("com.kazoup.srv.textanalyzer.Create", "id required")
	}

	if len(req.Index) == 0 {
		return errors.BadRequest("com.kazoup.srv.textanalyzer.Create", "index required")
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.textanalyzer.Create", err.Error())
	}

	go func() {
		if err := ta.Client.Publish(ctx, ta.Client.NewPublication(globals.ExtractEntitiesTopic, &enrich_proto.EnrichMessage{
			Index:  req.Index,
			Id:     req.Id,
			UserId: uID,
		})); err != nil {
			log.Println("ERROR publishing ExtractEntitiesTopic message", err)
		}
	}()

	return nil
}
