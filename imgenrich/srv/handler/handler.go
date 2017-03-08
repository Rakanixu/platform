package handler

import (
	proto "github.com/kazoup/platform/imgenrich/srv/proto/imgenrich"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type ImgEnrich struct {
	Client client.Client
}

func (ie *ImgEnrich) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("com.kazoup.srv.imgenrich.Create", "id required")
	}

	if len(req.Index) == 0 {
		return errors.BadRequest("com.kazoup.srv.imgenrich.Create", "index required")
	}

	if len(req.Type) == 0 {
		return errors.BadRequest("com.kazoup.srv.imgenrich.Create", "type required")
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.imgenrich.Create", err.Error())
	}

	go func() {
		if req.Type == globals.FileType {
			if err := ie.Client.Publish(ctx, ie.Client.NewPublication(globals.ImgEnrichTopic, &enrich_proto.EnrichMessage{
				Index:  req.Index,
				Id:     req.Id,
				UserId: uID,
			})); err != nil {
				log.Println("ERROR publishing ImgEnrichTopic message", err)
			}
		}

		if req.Type == globals.TypeDatasource {
			ids, err := retrieveImgFilesNotProcessed(ctx, ie.Client, req.Id, req.Index)
			if err != nil {
				log.Println("ERROR retireving image files not being process yet", err)
			}

			// Publish msg for all files not being process yet
			for _, v := range ids {
				if err := ie.Client.Publish(ctx, ie.Client.NewPublication(globals.ImgEnrichTopic, &enrich_proto.EnrichMessage{
					Index:  req.Index,
					Id:     v,
					UserId: uID,
				})); err != nil {
					log.Println("ERROR publishing ImgEnrichTopic message", err)
				}
			}
		}
	}()

	return nil
}
