package handler

import (
	"encoding/json"
	elastic "github.com/kazoup/platform/elastic/srv/proto/elastic"
	proto "github.com/kazoup/platform/policy/srv/proto/policy"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"time"
)

// Policy struct
type Policy struct{}

// Create srv handler
func (p *Policy) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Name) <= 0 {
		return errors.BadRequest("go.micro.srv.policy", "name required")
	}

	// Try read record, if exists, let consumer know policy name is in use
	srvReadReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Read",
		&elastic.ReadRequest{
			Index: "policies",
			Type:  "policy",
			Id:    req.Name,
		},
	)
	srvReadRsp := &elastic.ReadResponse{}

	if err := client.Call(ctx, srvReadReq, srvReadRsp); err == nil {
		return errors.BadRequest("go.micro.srv.policy", "Policy name already exists")
	}

	// Continue checking required fields exists
	if len(req.Filter) <= 0 {
		return errors.BadRequest("go.micro.srv.policy", "filter required")
	}

	if len(req.FilterRaw) <= 0 {
		return errors.BadRequest("go.micro.srv.policy", "filter_raw required")
	}

	// Fill in optional and record associate values
	req.Created = time.Now().String()

	if len(req.CreatedBy) <= 0 {
		req.CreatedBy = "unknown"
	}

	data, err := json.Marshal(req)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.policy", err.Error())
	}

	srvCreateReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Create",
		&elastic.CreateRequest{
			Index: "policies",
			Type:  "policy",
			Id:    req.Name,
			Data:  string(data),
		},
	)
	srvCreateRsp := &elastic.CreateResponse{}

	if err := client.Call(ctx, srvCreateReq, srvCreateRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.policy", err.Error())
	}

	return nil
}

// Read srv handler
func (p *Policy) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.Name) <= 0 {
		return errors.BadRequest("go.micro.srv.policy", "name required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Read",
		&elastic.ReadRequest{
			Index: "policies",
			Type:  "policy",
			Id:    req.Name,
		},
	)
	srvRsp := &elastic.ReadResponse{}

	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.policy", err.Error())
	}

	var result map[string]interface{}

	json.Unmarshal([]byte(srvRsp.Result), &result)

	// Those values are required, so never got nil interface
	rsp.Name = result["name"].(string)
	rsp.Created = result["created"].(string)
	rsp.Filter = result["filter"].(string)
	rsp.FilterRaw = result["filter_raw"].(string)

	// Optional values, set zero values for nil interfaces
	if result["created_by"] == nil {
		result["created_by"] = ""
	}
	rsp.CreatedBy = result["created_by"].(string)

	if result["is_archive_policy"] == nil {
		result["is_archive_policy"] = false
	}
	rsp.IsArchivePolicy = result["is_archive_policy"].(bool)

	if result["is_deletion_policy"] == nil {
		result["is_deletion_policy"] = false
	}
	rsp.IsDeletionPolicy = result["is_deletion_policy"].(bool)

	return nil
}

// Delete srv handler
func (p *Policy) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Name) <= 0 {
		return errors.BadRequest("go.micro.srv.policy", "name required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Delete",
		&elastic.DeleteRequest{
			Index: "policies",
			Type:  "policy",
			Id:    req.Name,
		},
	)
	srvRsp := &elastic.DeleteResponse{}

	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.policy", err.Error())
	}

	return nil
}

// List srv handler
func (p *Policy) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Search",
		&elastic.SearchRequest{
			Index:  "policies",
			Type:   "policy",
			Limit:  1000000,
			Offset: 0,
		},
	)
	srvRsp := &elastic.SearchResponse{}

	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.policy", err.Error())
	}

	var input map[string]interface{}
	var result []*proto.ReadResponse

	json.Unmarshal([]byte(srvRsp.Result), &input)
	// Iterate hits and generate slice of ReadResponse
	for _, hit := range input["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var isArchivePolicy, isDeletionPolicy bool

		if hit.(map[string]interface{})["_source"].(map[string]interface{})["is_archive_policy"] == nil {
			isArchivePolicy = false
		} else {
			isArchivePolicy = true
		}

		if hit.(map[string]interface{})["_source"].(map[string]interface{})["is_deletion_policy"] == nil {
			isDeletionPolicy = false
		} else {
			isDeletionPolicy = true
		}

		result = append(result, &proto.ReadResponse{
			Name:             hit.(map[string]interface{})["_source"].(map[string]interface{})["name"].(string),
			Created:          hit.(map[string]interface{})["_source"].(map[string]interface{})["created"].(string),
			CreatedBy:        hit.(map[string]interface{})["_source"].(map[string]interface{})["created_by"].(string),
			Filter:           hit.(map[string]interface{})["_source"].(map[string]interface{})["filter"].(string),
			FilterRaw:        hit.(map[string]interface{})["_source"].(map[string]interface{})["filter_raw"].(string),
			IsArchivePolicy:  isArchivePolicy,
			IsDeletionPolicy: isDeletionPolicy,
		})
	}

	rsp.Result = result

	return nil
}
