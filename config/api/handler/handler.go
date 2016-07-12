package handler

import (
	"encoding/json"
	"net/http"

	config "github.com/kazoup/platform/config/srv/proto/config"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

// Config struct
type Config struct{}

// Status handler
func (c *Config) Status(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.config",
		"Config.Status",
		&config.StatusRequest{},
	)
	srvRes := &config.StatusResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	b, err := json.Marshal(srvRes)
	if err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	rsp.Body = string(b)
	rsp.StatusCode = http.StatusOK

	return nil
}

// SetElasticSettings API handler, sets ElasticSearch settings (check es_settings.json in repo) for files index
func (c *Config) SetElasticSettings(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.config",
		"Config.SetElasticSettings",
		&config.SetElasticSettingsRequest{},
	)
	srvRes := &config.SetElasticSettingsResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// SetElasticMapping API handler, sets ElasticSearch mapping for files index and file documents.(Check es_mapping_files.json in this repo)
func (c *Config) SetElasticMapping(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.config",
		"Config.SetElasticMapping",
		&config.SetElasticMappingRequest{},
	)
	srvRes := &config.SetElasticMappingResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// SetFlags API handler, posts appliance flags into ElasticSearch.(Check es_flags.json in this repo)
func (c *Config) SetFlags(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.config",
		"Config.SetFlags",
		&config.SetFlagsRequest{},
	)
	srvRes := &config.SetFlagsResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}
