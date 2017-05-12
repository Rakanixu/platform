package handler

import (
	"errors"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"golang.org/x/net/context"
)

const (
	QUOTA_EXCEEDED_MSG = "Quota for Speech to text service exceeded."
)

type Service struct {
	quota quota.Checker
}

func (s *Service) EnrichFile(ctx context.Context, req *proto_audio.EnrichFileRequest, rsp *proto_audio.EnrichFileResponse) error {
	if err := validate.Exists(req.Id, req.Index); err != nil {
		return err
	}

	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	_, _, rate, _, quota, ok := quota.Check(ctx, globals.AUDIO_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.AUDIO_SERVICE_NAME, "EnrichFile", "", errors.New("quota.Check"))
	}

	// Quota exceded, respond sync and do not initiate go routines
	if rate-quota > 0 {
		rsp.Info = QUOTA_EXCEEDED_MSG
		return nil
	}

	return nil
}

func (s *Service) EnrichDatasource(ctx context.Context, req *proto_audio.EnrichDatasourceRequest, rsp *proto_audio.EnrichDatasourceResponse) error {
	if err := validate.Exists(req.Id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_audio.HealthRequest, rsp *proto_audio.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
