package handler

import (
	"github.com/kazoup/platform/entities/srv/proto/entities"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	proto_quota "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type Service struct{}

// ExtractFile handler
func (s *Service) ExtractFile(ctx context.Context, req *proto_entities.ExtractFileRequest, rsp *proto_entities.ExtractFileResponse) error {
	if err := validate.Exists(ctx, req.Id, req.Index); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

	// Check Quota first
	qreq := srv.Client().NewRequest(
		globals.QUOTA_SERVICE_NAME,
		"Quota.Read",
		&proto_quota.ReadRequest{
			Srv: globals.ENTITIES_SERVICE_NAME,
		},
	)
	qrsp := &proto_quota.ReadResponse{}
	if err := srv.Client().Call(ctx, qreq, qrsp); err != nil {
		log.Println("Error calling Quota.Read", err)
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Entity extraction service exceeded."
		return nil
	}

	return nil
}

// Health service handler
func (s *Service) Health(ctx context.Context, req *proto_entities.HealthRequest, rsp *proto_entities.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
