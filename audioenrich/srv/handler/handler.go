package handler

import (
	proto "github.com/kazoup/platform/audioenrich/srv/proto/audioenrich"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type Enrich struct {
	Client client.Client
}

func (e *Enrich) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("com.kazoup.srv.audioenrich.Create", "id required")
	}

	if len(req.Index) == 0 {
		return errors.BadRequest("com.kazoup.srv.audioenrich.Create", "index required")
	}

	if len(req.Type) == 0 {
		return errors.BadRequest("com.kazoup.srv.audioenrich.Create", "type required")
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.audioenrich.Create", err.Error())
	}

	go func() {
		if req.Type == globals.FileType {
			if err := e.Client.Publish(ctx, e.Client.NewPublication(globals.AudioEnrichTopic, &enrich_proto.EnrichMessage{
				Index:  req.Index,
				Id:     req.Id,
				UserId: uID,
			})); err != nil {
				log.Println("ERROR publishing AudioEnrichTopic message", err)
			}
		}

		if req.Type == globals.TypeDatasource {
			ids, err := retrieveFilesNotProcessed(ctx, req.Id, req.Index)
			if err != nil {
				log.Println("ERROR retireving audio files not being process yet", err)
			}

			// Publish msg for all files not being process yet
			for _, v := range ids {
				if err := e.Client.Publish(ctx, e.Client.NewPublication(globals.AudioEnrichTopic, &enrich_proto.EnrichMessage{
					Index:  req.Index,
					Id:     v,
					UserId: uID,
				})); err != nil {
					log.Println("ERROR publishing AudioEnrichTopic message", err)
				}
			}
		}
	}()

	return nil
}
