package handler

import (
	"github.com/kazoup/platform/image/srv/proto/image"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	proto_quota "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type Service struct{}

func (s *Service) EnrichFile(ctx context.Context, req *proto_image.EnrichFileRequest, rsp *proto_image.EnrichFileResponse) error {
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
			Srv: globals.IMAGE_SERVICE_NAME,
		},
	)
	qrsp := &proto_quota.ReadResponse{}
	if err := srv.Client().Call(ctx, qreq, qrsp); err != nil {
		log.Println("Error calling Quota.Read", err)
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Image content service exceeded."
		return nil
	}

	return nil
}

func (s *Service) EnrichDatasource(ctx context.Context, req *proto_image.EnrichDatasourceRequest, rsp *proto_image.EnrichDatasourceResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_image.HealthRequest, rsp *proto_image.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
