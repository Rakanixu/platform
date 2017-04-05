package handler

import (
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	quota_proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/kazoup/platform/sentiment/srv/proto/sentiment"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type Service struct{}

// AnalyzeFile handler
func (s *Service) AnalyzeFile(ctx context.Context, req *proto_sentiment.AnalyzeFileRequest, rsp *proto_sentiment.AnalyzeFileResponse) error {
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
		&quota_proto.ReadRequest{
			Srv: globals.SENTIMENT_SERVICE_NAME,
		},
	)
	qrsp := &quota_proto.ReadResponse{}
	if err := srv.Client().Call(ctx, qreq, qrsp); err != nil {
		return err
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Entity extraction service exceeded."
		return nil
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_sentiment.HealthRequest, rsp *proto_sentiment.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
