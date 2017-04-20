package handler

import (
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/validate"
	"github.com/kazoup/platform/user/srv/proto/user"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Service struct{}

// Create File handler
func (s *Service) Read(ctx context.Context, req *proto_user.ReadRequest, rsp *proto_user.ReadResponse) error {
	if err := validate.Exists(ctx, req.Index, req.Id); err != nil {
		return err
	}

	res, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: req.Index,
		Type:  globals.UserType,
		Id:    globals.GetMD5Hash(req.Id),
	})

	if err != nil {
		return errors.InternalServerError(globals.USER_SERVICE_NAME, err.Error())
	}

	rsp.Result = res.Result

	return nil
}

func (s *Service) Health(ctx context.Context, req *proto_user.HealthRequest, rsp *proto_user.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
