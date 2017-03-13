package handler

import (
	slib "github.com/kazoup/platform/lib/stripe"
	proto "github.com/kazoup/platform/profile/srv/proto/profile"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Profile struct
type Profile struct {
	Client client.Client
}

// Read Profile handler
func (p *Profile) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.StripeId) == 0 {
		return errors.InternalServerError("com.kazoup.srv.profile.Read", "stripe_id required")
	}

	m, pf, err := slib.GetCustomer(req.StripeId)
	if err != nil {
		return err
	}

	rsp.Profile = &proto.Profile{
		Email:   m,
		Profile: pf,
	}

	return nil
}

// Health
func (p *Profile) Health(ctx context.Context, req *proto.HealthRequest, rsp *proto.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
