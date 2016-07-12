package handler

import (
	"encoding/json"
	"net/http"

	proto "github.com/kazoup/platform/policy/srv/proto/policy"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

type Policy struct{}

// Create API handler
func (p *Policy) Create(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}

	// Unmarshal unknown JSON
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Required
	if input["name"] == nil {
		return errors.BadRequest("go.micro.api.policy", "name required")
	}

	if input["filter"] == nil {
		return errors.BadRequest("go.micro.api.policy", "filter required")
	}

	if input["filter_raw"] == nil {
		return errors.BadRequest("go.micro.api.policy", "filter_raw required")
	}

	if input["is_archive_policy"] == nil {
		return errors.BadRequest("go.micro.api.policy", "is_archive_policy required")
	}

	if input["is_deletion_policy"] == nil {
		return errors.BadRequest("go.micro.api.policy", "is_deletion_policy required")
	}

	// Optional
	if input["created_by"] == nil {
		input["created_by"] = ""
	}

	// Marshal unknown JSON (filter)
	filter, err := json.Marshal(input["filter"])
	if err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	srvReq := client.NewRequest(
		"go.micro.srv.policy",
		"Policy.Create",
		&proto.CreateRequest{
			Name:             input["name"].(string),
			Filter:           string(filter),
			FilterRaw:        input["filter_raw"].(string),
			IsArchivePolicy:  input["is_archive_policy"].(bool),
			IsDeletionPolicy: input["is_archive_policy"].(bool),
			Created:          "",
			CreatedBy:        input["created_by"].(string),
		},
	)
	srvRsp := &proto.CreateResponse{}

	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// Read API handler
func (p *Policy) Read(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}

	// Unmarshal unknown JSON
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	if input["name"] == nil {
		return errors.BadRequest("go.micro.api.policy", "name required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.policy",
		"Policy.Read",
		&proto.ReadRequest{
			Name: input["name"].(string),
		},
	)
	srvRsp := &proto.ReadResponse{}

	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Helper to manage JSON response
	policy, err := PolicyToStringJSON(srvRsp)
	if err != nil {
		return err
	}

	rsp.Body = string(policy)
	rsp.StatusCode = http.StatusOK

	return nil
}

// Delete API handler
func (p *Policy) Delete(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}

	// Unmarshal unknown JSON
	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Required
	if input["name"] == nil {
		return errors.BadRequest("go.micro.api.policy", "name required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.policy",
		"Policy.Delete",
		&proto.DeleteRequest{
			Name: input["name"].(string),
		},
	)
	srvRsp := &proto.DeleteResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	rsp.Body = `{}`
	rsp.StatusCode = http.StatusOK

	return nil
}

// List API handler
func (p *Policy) List(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.policy",
		"Policy.List",
		&proto.ListRequest{},
	)
	srvRsp := &proto.ListResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.api.policy", err.Error())
	}

	// Helper to manage JSON response
	b, err := PoliciesToStringJSON(srvRsp.Result)
	if err != nil {
		return err
	}

	rsp.Body = string(b)

	rsp.StatusCode = http.StatusOK

	return nil
}
