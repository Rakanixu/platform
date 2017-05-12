package handler

import (
	"errors"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"github.com/kazoup/platform/sentiment/srv/proto/sentiment"
	"golang.org/x/net/context"
)

const (
	QUOTA_EXCEEDED_MSG = "Quota for Entity extraction service exceeded."
)

type Service struct{}

// AnalyzeFile handler
func (s *Service) AnalyzeFile(ctx context.Context, req *proto_sentiment.AnalyzeFileRequest, rsp *proto_sentiment.AnalyzeFileResponse) error {
	if err := validate.Exists(req.Id, req.Index); err != nil {
		return err
	}

	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	// Check Quota first
	_, _, rate, _, quota, ok := quota.Check(ctx, globals.SENTIMENT_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.SENTIMENT_SERVICE_NAME, "AnalyzeFile", "", errors.New("quota.Check"))
	}

	// Quota exceded, respond sync and do not initiate go routines
	if rate-quota > 0 {
		rsp.Info = QUOTA_EXCEEDED_MSG
		return nil
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_sentiment.HealthRequest, rsp *proto_sentiment.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
