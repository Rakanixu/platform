package handler

import (
	proto "github.com/kazoup/platform/audioenrich/srv/proto/audioenrich"
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	quota_proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type AudioEnrich struct {
	Client client.Client
}

func (ae *AudioEnrich) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
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

	// We only want to respond now if the endpoint was hitted for specific file
	if req.Type == globals.FileType {
		// Check Quota first
		qreq := ae.Client.NewRequest(
			globals.QUOTA_SERVICE_NAME,
			"Quota.Read",
			&quota_proto.ReadRequest{
				Srv: globals.AUDIOENRICH_SERVICE_NAME,
			},
		)
		qrsp := &quota_proto.ReadResponse{}
		if err := ae.Client.Call(ctx, qreq, qrsp); err != nil {
			log.Println("Error calling Quota.Read", err)
		}

		// Quota exceded, respond sync and do not initiate go routines
		if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
			rsp.Info = "Quota for Speech to text service exceeded."
			return nil
		}
	}

	go func() {
		if req.Type == globals.FileType {
			if err := ae.Client.Publish(ctx, ae.Client.NewPublication(globals.AudioEnrichTopic, &enrich_proto.EnrichMessage{
				Index:  req.Index,
				Id:     req.Id,
				UserId: uID,
				Notify: true,
			})); err != nil {
				log.Println("ERROR publishing AudioEnrichTopic message", err)
			}
		}

		if req.Type == globals.TypeDatasource {
			ids, err := retrieveAudioFilesNotProcessed(ctx, ae.Client, req.Id, req.Index)
			if err != nil {
				log.Println("ERROR retireving audio files not being process yet", err)
			}

			// Publish msg for all files not being process yet
			for _, v := range ids {
				if err := ae.Client.Publish(ctx, ae.Client.NewPublication(globals.AudioEnrichTopic, &enrich_proto.EnrichMessage{
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
