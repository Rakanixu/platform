package handler

import (
	"errors"
	"github.com/kazoup/platform/agent/srv/proto/agent"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/quota"
	"github.com/kazoup/platform/lib/utils"
	"github.com/kazoup/platform/lib/validate"
	"golang.org/x/net/context"
)

const (
	QUOTA_EXCEEDED_MSG = "Quota for Agent service exceeded."
)

// Service structure
type Service struct{}

// Saves Kazup file
func (s *Service) Save(ctx context.Context, req *proto_agent.SaveRequest, rsp *proto_agent.SaveResponse) error {
	// Check if data parameter exists in the request
	if err := validate.Exists(req.Data); err != nil {
		return err
	}

	// Extract user id
	uID, err := utils.ParseUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	// Get quota limit and current rate
	_, _, rate, _, quota, ok := quota.Check(ctx, globals.AGENT_SERVICE_NAME, uID)
	if !ok {
		return platform_errors.NewPlatformError(globals.AGENT_SERVICE_NAME, "Save", "", 403, errors.New("quota.Check"))
	}

	// Quota exceded, stop execution
	if rate-quota > 0 {
		rsp.Info = QUOTA_EXCEEDED_MSG
		return nil
	}

	return nil
}
