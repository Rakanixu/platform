package handler

import (
	"encoding/json"
	proto "github.com/kazoup/platform/config/srv/proto/config"
	elastic "github.com/kazoup/platform/elastic/srv/proto/elastic"
	flag "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Config struct
type Config struct {
	Client             client.Client
	ElasticServiceName string
	ESSettings         *[]byte
	ESFlags            *[]byte
	ESMapping          *[]byte
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

// SetElasticSettings handler, sets ElasticSearch settings (check es_settings.json in repo) for files index
func (c *Config) SetElasticSettings(ctx context.Context, req *proto.SetElasticSettingsRequest, rsp *proto.SetElasticSettingsResponse) error {
	srvReq := c.Client.NewRequest(
		c.ElasticServiceName,
		"Elastic.CreateIndexWithSettings",
		&elastic.CreateIndexWithSettingsRequest{
			Index:    "files",
			Settings: string(*c.ESSettings),
		},
	)
	srvRes := &elastic.CreateIndexWithSettingsResponse{}

	if err := c.Client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	return nil
}

// SetElasticSettings handler, sets ElasticSearch mapping for files index and file documents.(Check es_mapping_files.json in this repo)
func (c *Config) SetElasticMapping(ctx context.Context, req *proto.SetElasticMappingRequest, rsp *proto.SetElasticMappingResponse) error {

	srvReq := c.Client.NewRequest(
		c.ElasticServiceName,
		"Elastic.PutMappingFromJSON",
		&elastic.PutMappingFromJSONRequest{
			Index:   "files",
			Type:    "file",
			Mapping: string(*c.ESMapping),
		},
	)
	srvRes := &elastic.PutMappingFromJSONResponse{}

	if err := c.Client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.config", err.Error())
	}

	return nil
}

// SetFlags handler, post kazoup appliance flags into elastic search. (Check es_flags.json)
func (c *Config) SetFlags(ctx context.Context, req *proto.SetFlagsRequest, rsp *proto.SetFlagsResponse) error {
	var flagsSlice []interface{}

	flags := *c.ESFlags
	if err := json.Unmarshal(flags, &flagsSlice); err != nil {
		return errors.InternalServerError("go.micro.srv.config", err.Error())
	}

	for _, v := range flagsSlice {
		srvReq := c.Client.NewRequest(
			c.ElasticServiceName,
			"Flag.Create",
			// We ensure es_flags.json contains proper JSON, so fields will be there always
			&flag.CreateRequest{
				Key:         v.(map[string]interface{})["key"].(string),
				Description: v.(map[string]interface{})["description"].(string),
				Value:       v.(map[string]interface{})["value"].(bool),
			},
		)
		srvRes := &elastic.CreateResponse{}

		if err := c.Client.Call(ctx, srvReq, srvRes); err != nil {
			return errors.InternalServerError("go.micro.api.config", err.Error())
		}
	}

	return nil
}
