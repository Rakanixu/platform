package handler

import (
	"encoding/json"
	data "github.com/kazoup/platform/config/srv/data"
	proto "github.com/kazoup/platform/config/srv/proto/config"
	flag "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"log"
)

const (
	Desktop    = "desktop"
	Enterprise = "enterprise"
)

// Config struct
type Config struct {
	Client        client.Client
	DbServiceName string
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

// SetFlags handler, post kazoup appliance flags into elastic search. (Check es_flags.json)
func (c *Config) SetFlags(ctx context.Context, req *proto.SetFlagsRequest, rsp *proto.SetFlagsResponse) error {
	var asset string

	if len(req.Type) == 0 {
		return errors.InternalServerError("go.micro.srv.config", "type required")
	}

	if req.Type == Desktop {
		asset = "data/es_desktop_flags.json"
	}

	if req.Type == Enterprise {
		asset = "data/es_flags.json"
	}

	es_flags, err := data.Asset(asset)
	if err != nil {
		// Asset was not found.
		log.Fatal(err)
	}

	var flagsSlice []interface{}

	if err := json.Unmarshal(es_flags, &flagsSlice); err != nil {
		return errors.InternalServerError("go.micro.srv.config", err.Error())
	}

	for _, v := range flagsSlice {
		srvReq := c.Client.NewRequest(
			c.DbServiceName,
			"Flag.Create",
			&flag.CreateRequest{
				Key:         v.(map[string]interface{})["key"].(string),
				Description: v.(map[string]interface{})["description"].(string),
				Value:       v.(map[string]interface{})["value"].(bool),
			},
		)
		srvRes := &flag.CreateResponse{}

		if err := c.Client.Call(ctx, srvReq, srvRes); err != nil {
			return errors.InternalServerError("go.micro.api.config", err.Error())
		}
	}

	return nil
}
