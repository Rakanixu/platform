package handler

import (
	"encoding/json"
	"fmt"
	"strconv"

	elasticsearch "github.com/kazoup/platform/elastic/srv/proto/elastic"
	proto "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

// Flag struct
type Flag struct{}

// Create srv handler
func (f *Flag) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Key) <= 0 || len(req.Description) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Create", "Fields required")
	}

	data, _ := json.Marshal(req)

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Create",
		&elasticsearch.CreateRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Id for flags will be the key given, so we can RUD easily
			Data:  string(data),
		},
	)
	srvRsp := &elasticsearch.CreateResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Create", err.Error())
	}

	return nil
}

// Read srv handler
func (f *Flag) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.Key) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Read", "Flag key required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Read",
		&elasticsearch.ReadRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvRsp := &elasticsearch.ReadResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Read", err.Error())
	}

	// micro service returns a string, let's map it and return
	var input map[string]interface{}

	json.Unmarshal([]byte(srvRsp.Result), &input)

	rsp.Key = fmt.Sprintf("%v", input["key"])
	rsp.Description = fmt.Sprintf("%v", input["description"])
	rsp.Value, _ = strconv.ParseBool(fmt.Sprintf("%v", input["value"]))

	return nil
}

// Flip srv handler
func (f *Flag) Flip(ctx context.Context, req *proto.FlipRequest, rsp *proto.FlipResponse) error {
	if len(req.Key) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Flip", "Flag key required")
	}

	// Read the record to flip
	srvReadReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Read",
		&elasticsearch.ReadRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvReadRsp := &elasticsearch.ReadResponse{}
	if err := client.Call(ctx, srvReadReq, srvReadRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Flip", err.Error())
	}

	// Flip value
	var input map[string]interface{}

	json.Unmarshal([]byte(srvReadRsp.Result), &input)
	value, _ := strconv.ParseBool(fmt.Sprintf("%v", input["value"]))
	input["value"] = !value
	data, _ := json.Marshal(input)

	// Update the record
	srvUpdateReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Update",
		&elasticsearch.UpdateRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
			Data:  string(data),
		},
	)
	srvUpdateRsp := &elasticsearch.UpdateResponse{}
	if err := client.Call(ctx, srvUpdateReq, srvUpdateRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Flip", err.Error())
	}

	return nil
}

// Delete srv handler
func (f *Flag) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Key) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Delete", "Flag key required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Delete",
		&elasticsearch.DeleteRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvRsp := &elasticsearch.DeleteResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Delete", err.Error())
	}

	return nil
}

// List srv handler
func (f *Flag) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Search in flags index
	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Search",
		&elasticsearch.SearchRequest{
			Index:  "flags", // Hardcoded index for flags
			Type:   "flag",  // Hardcoded type ...
			Query:  "*",     // No filter
			Offset: 0,       // From the first one
			Limit:  1000000, // Hope you've got less than a million flags ...
		},
	)
	srvRsp := &elasticsearch.SearchResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.List", err.Error())
	}

	var input map[string]interface{}
	var result []*proto.ReadResponse

	json.Unmarshal([]byte(srvRsp.Result), &input)
	// Iterate hits and generate slice of ReadResponse
	for _, hit := range input["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var val bool

		if hit.(map[string]interface{})["_source"].(map[string]interface{})["value"] == nil {
			val = false
		} else {
			val = true
		}

		result = append(result, &proto.ReadResponse{
			Key:         hit.(map[string]interface{})["_source"].(map[string]interface{})["key"].(string),
			Description: hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string),
			Value:       val,
		})
	}

	rsp.Result = result

	return nil
}
