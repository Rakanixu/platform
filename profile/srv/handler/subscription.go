package handler

import (
	"github.com/kazoup/platform/lib/globals"
	slib "github.com/kazoup/platform/lib/stripe"
	proto "github.com/kazoup/platform/profile/srv/proto/subscription"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Subscription struct
type Subscription struct {
	Client client.Client
}

// Update subscription handler
func (s *Subscription) Update(ctx context.Context, req *proto.UpdateSubRequest, rsp *proto.UpdateSubResponse) error {
	if len(req.StripeId) == 0 {
		return errors.InternalServerError("com.kazoup.srv.profile.Create", "stripe_id required")
	}

	if len(req.SubId) == 0 {
		return errors.InternalServerError("com.kazoup.srv.profile.Create", "sub_id required")
	}

	var pt string
	// Upgrade
	if req.UpgradeSubscription {
		if len(req.CheckoutTokenId) == 0 {
			return errors.InternalServerError("com.kazoup.srv.profile.Create", "checkout_token_id required")
		}

		// Save credit card associated with user
		if err := slib.SaveCard(req.StripeId, req.CheckoutTokenId); err != nil {
			return errors.InternalServerError("com.kazoup.srv.profile.Create", err.Error())
		}
		pt = globals.PRODUCT_TYPE_TEAM
	}

	// Downgrade
	if !req.UpgradeSubscription {
		pt = globals.PRODUCT_TYPE_PERSONAL
	}

	// Update to team for now
	if err := slib.UpdateSubscription(req.SubId, pt); err != nil {
		return errors.InternalServerError("com.kazoup.srv.profile.Create", err.Error())
	}

	return nil
}
