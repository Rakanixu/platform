package handler

import (
	proto "github.com/kazoup/platform/docenrich/srv/proto/docenrich"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type DocEnrich struct {
	Client client.Client
}

func (de *DocEnrich) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("com.kazoup.srv.docenrich.Create", "id required")
	}

	if len(req.Index) == 0 {
		return errors.BadRequest("com.kazoup.srv.docenrich.Create", "index required")
	}

	if len(req.Type) == 0 {
		return errors.BadRequest("com.kazoup.srv.docenrich.Create", "type required")
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.docenrich.Create", err.Error())
	}

	go func() {
		if req.Type == globals.FileType {
			if err := de.Client.Publish(ctx, de.Client.NewPublication(globals.DocEnrichTopic, &enrich_proto.EnrichMessage{
				Index:  req.Index,
				Id:     req.Id,
				UserId: uID,
				Notify: true,
			})); err != nil {
				log.Println("ERROR publishing DocEnrichTopic message", err)
			}
		}

		if req.Type == globals.TypeDatasource {
			ids, err := retrieveDocFilesNotProcessed(ctx, de.Client, req.Id, req.Index)
			if err != nil {
				log.Println("ERROR retireving doc files not being process yet", err)
			}

			// Publish msg for all files not being process yet
			for _, v := range ids {
				if err := de.Client.Publish(ctx, de.Client.NewPublication(globals.DocEnrichTopic, &enrich_proto.EnrichMessage{
					Index:  req.Index,
					Id:     v,
					UserId: uID,
				})); err != nil {
					log.Println("ERROR publishing DocEnrichTopic message", err)
				}
			}
		}
	}()

	return nil
}
