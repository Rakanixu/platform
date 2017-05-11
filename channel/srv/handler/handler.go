package handler

import (
	"github.com/kazoup/platform/channel/srv/proto/channel"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Service struct{}

// Create File handler
func (s *Service) Read(ctx context.Context, req *proto_channel.ReadRequest, rsp *proto_channel.ReadResponse) error {
	if err := validate.Exists(ctx, req.Index, req.Id); err != nil {
		return err
	}

	res, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: req.Index,
		Type:  globals.ChannelType,
		Id:    utils.GetMD5Hash(req.Id),
	})

	if err != nil {
		return errors.InternalServerError(globals.CHANNEL_SERVICE_NAME, err.Error())
	}

	rsp.Result = res.Result

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_channel.HealthRequest, rsp *proto_channel.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
