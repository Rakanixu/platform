package handler

import (
	"github.com/kazoup/platform/document/srv/proto/document"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	proto_quota "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type Service struct{}

func (de *Service) EnrichFile(ctx context.Context, req *proto_document.EnrichFileRequest, rsp *proto_document.EnrichFileResponse) error {
	if err := validate.Exists(ctx, req.Id, req.Index); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return platform_errors.ErrInvalidCtx
	}

	// Check Quota first
	qreq := srv.Client().NewRequest(
		globals.QUOTA_SERVICE_NAME,
		"Quota.Read",
		&proto_quota.ReadRequest{
			Srv: globals.DOCUMENT_SERVICE_NAME,
		},
	)
	qrsp := &proto_quota.ReadResponse{}
	if err := srv.Client().Call(ctx, qreq, qrsp); err != nil {
		return err
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Document content extraction service exceeded."
		return nil
	}

	return nil
}

func (s *Service) EnrichDatasource(ctx context.Context, req *proto_document.EnrichDatasourceRequest, rsp *proto_document.EnrichDatasourceResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_document.HealthRequest, rsp *proto_document.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
