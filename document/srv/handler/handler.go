package handler

import (
	"errors"
	"github.com/kazoup/platform/document/srv/proto/document"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"golang.org/x/net/context"
)

const (
	QUOTA_EXCEEDED_MSG = "Quota for Document content extraction service exceeded."
)

type Service struct{}

func (de *Service) EnrichFile(ctx context.Context, req *proto_document.EnrichFileRequest, rsp *proto_document.EnrichFileResponse) error {
	if err := validate.Exists(ctx, req.Id, req.Index); err != nil {
		return err
	}

	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	// Check Quota first
	_, _, rate, _, quota, ok := quota.Check(ctx, globals.DOCUMENT_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.DOCUMENT_SERVICE_NAME, "EnrichFile", "", errors.New("quota.Check"))
	}

	// Quota exceded, respond sync and do not initiate go routines
	if rate-quota > 0 {
		rsp.Info = QUOTA_EXCEEDED_MSG
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
