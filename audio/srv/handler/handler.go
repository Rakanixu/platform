package handler

import (
	"github.com/kazoup/platform/audio/srv/proto/audio"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	quota_proto "github.com/kazoup/platform/quota/srv/proto/quota"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Service struct{}

func (s *Service) EnrichFile(ctx context.Context, req *proto_audio.EnrichFileRequest, rsp *proto_audio.EnrichFileResponse) error {
	if err := validate.Exists(ctx, req.Id, req.Index); err != nil {
		return err
	}

	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.New("Cant get srv from context", "", 500)
	}

	// Check Quota for audio service
	qreq := srv.Client().NewRequest(
		globals.QUOTA_SERVICE_NAME,
		"Quota.Read",
		&quota_proto.ReadRequest{
			Srv: globals.AUDIO_SERVICE_NAME,
		},
	)
	qrsp := &quota_proto.ReadResponse{}
	if err := srv.Client().Call(ctx, qreq, qrsp); err != nil {
		return err
	}

	// Quota exceded, respond sync and do not initiate go routines
	if qrsp.Quota.Rate-qrsp.Quota.Quota > 0 {
		rsp.Info = "Quota for Speech to text service exceeded."
		return nil
	}

	return nil
}

func (s *Service) EnrichDatasource(ctx context.Context, req *proto_audio.EnrichDatasourceRequest, rsp *proto_audio.EnrichDatasourceResponse) error {
	if err := validate.Exists(ctx, req.Id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_audio.HealthRequest, rsp *proto_audio.HealthResponse) error {
	rsp.Info = "OK"
	rsp.Status = 200

	return nil
}
