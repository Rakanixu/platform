package handler

import (
	"errors"
	"github.com/kazoup/platform/entities/srv/proto/entities"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/validate"
	"golang.org/x/net/context"
)

type Service struct{}

// ExtractFile handler
func (s *Service) ExtractFile(ctx context.Context, req *proto_entities.ExtractFileRequest, rsp *proto_entities.ExtractFileResponse) error {
	if err := validate.Exists(ctx, req.Id, req.Index); err != nil {
		return err
	}

	uID, err := globals.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	// Check Quota first
	_, _, rate, _, quota, ok := quota.Check(ctx, globals.ENTITIES_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.ENTITIES_SERVICE_NAME, "ExtractFile", "", errors.New("quota.Check"))
	}

	// Quota exceded, respond sync and do not initiate go routines
	if rate-quota > 0 {
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
