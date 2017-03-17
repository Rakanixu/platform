package handler

import (
	"github.com/kazoup/platform/lib/globals"
	enrich_proto "github.com/kazoup/platform/lib/protomsg"
	quota_proto "github.com/kazoup/platform/quota/srv/proto/quota"
	proto "github.com/kazoup/platform/sentimentanalyzer/srv/proto/sentimentanalyzer"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

type SentimentAnalyzer struct {
	Client client.Client
}

// Create handler publish a SINGLE FILE ExtractEntitiesTopic msg (in difference of imgenrich, docenrich, this endpoint DOES NOT SUPPORT DATASOURCE)
func (sa *SentimentAnalyzer) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("com.kazoup.srv.sentimentanalyzer.Create", "id required")
	}

	if len(req.Index) == 0 {
		return errors.BadRequest("com.kazoup.srv.sentimentanalyzer.Create", "index required")
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return errors.InternalServerError("com.kazoup.srv.sentimentanalyzer.Create", err.Error())
	}

	// Check Quota first
	qreq := sa.Client.NewRequest(
		globals.QUOTA_SERVICE_NAME,
		"Quota.Read",
		&quota_proto.ReadRequest{
			Srv: globals.SENTIMENTANALYZER_SERVICE_NAME,
		},
	)
	qrsp := &quota_proto.ReadResponse{}
	if err := sa.Client.Call(ctx, qreq, qrsp); err != nil {
		log.Println("Error calling Quota.Read", err)
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Entity extraction service exceeded."
		return nil
	}

	go func() {
		if err := sa.Client.Publish(ctx, sa.Client.NewPublication(globals.ExtractEntitiesTopic, &enrich_proto.EnrichMessage{
			Index:  req.Index,
			Id:     req.Id,
			UserId: uID,
			Notify: true,
		})); err != nil {
			log.Println("ERROR publishing ExtractEntitiesTopic message", err)
		}
	}()

	return nil
}
