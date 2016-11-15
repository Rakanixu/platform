package handler

import (
	proto "github.com/kazoup/platform/config/srv/proto/config"

	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

const (
	Desktop    = "desktop"
	Enterprise = "enterprise"
)

// Config struct
type Config struct {
	Client client.Client
}

// Status handler, retrieve kazoup appliance status
func (c *Config) Status(ctx context.Context, req *proto.StatusRequest, rsp *proto.StatusResponse) error {
	//TODO: implementation
	// mock response
	rsp.APPLIANCE_IS_CONFIGURED = false
	rsp.APPLIANCE_IS_DEMO = false
	rsp.APPLIANCE_IS_REGISTERED = true
	rsp.GIT_COMMIT_STRING = "asdfasdfasdfasdfasdfsdafhash"
	rsp.SMB_USER_EXISTS = true

	return nil
}
