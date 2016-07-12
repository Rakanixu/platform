package handler

import (
	"encoding/json"
	"net/http"

	proto "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

// Flag struct
type Flag struct{}

// Create API handler
func (f *Flag) Create(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var err error

	if len(req.Body) <= 0 {
		return errors.BadRequest("go.micro.api.flag", "Flag required")
	}

	var flag *proto.CreateRequest
	if err = json.Unmarshal([]byte(req.Body), &flag); err != nil {
		return errors.BadRequest("go.micro.api.flag", "Error parsing flag")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.flag",
		"Flag.Create",
		flag,
	)
	srvRsp := &proto.CreateResponse{}
	if err = client.Call(ctx, srvReq, srvRsp); err != nil {
		return err
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// Read API handler
func (f *Flag) Read(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var err error

	if len(req.Body) <= 0 {
		return errors.BadRequest("go.micro.api.flag", "Flag key required")
	}

	var key *proto.ReadRequest
	if err = json.Unmarshal([]byte(req.Body), &key); err != nil {
		return errors.BadRequest("go.micro.api.flag", "Error parsing key")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.flag",
		"Flag.Read",
		key,
	)
	srvRsp := &proto.ReadResponse{}
	if err = client.Call(ctx, srvReq, srvRsp); err != nil {
		return err
	}

	b, _ := json.Marshal(srvRsp)

	rsp.StatusCode = http.StatusOK
	rsp.Body = string(b)

	return nil
}

// Flip API handler
func (f *Flag) Flip(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var err error

	if len(req.Body) <= 0 {
		return errors.BadRequest("go.micro.api.flag", "Flag key required")
	}

	var key *proto.FlipRequest
	if err = json.Unmarshal([]byte(req.Body), &key); err != nil {
		return errors.BadRequest("go.micro.api.flag", "Error parsing key")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.flag",
		"Flag.Flip",
		key,
	)
	srvRsp := &proto.FlipResponse{}
	if err = client.Call(ctx, srvReq, srvRsp); err != nil {
		return err
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// Delete API handler
func (f *Flag) Delete(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var err error

	if len(req.Body) <= 0 {
		return errors.BadRequest("go.micro.api.flag", "Flag key required")
	}

	var key *proto.DeleteRequest
	if err = json.Unmarshal([]byte(req.Body), &key); err != nil {
		return errors.BadRequest("go.micro.api.flag", "Error parsing key")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.flag",
		"Flag.Delete",
		key,
	)
	srvRsp := &proto.FlipResponse{}
	if err = client.Call(ctx, srvReq, srvRsp); err != nil {
		return err
	}

	rsp.StatusCode = http.StatusOK
	rsp.Body = `{}`

	return nil
}

// List API handler
func (f *Flag) List(ctx context.Context, req *api.Request, rsp *api.Response) error {
	srvReq := client.NewRequest(
		"go.micro.srv.flag",
		"Flag.List",
		&proto.ListRequest{},
	)
	srvRsp := &proto.ListResponse{}
	if err := client.Call(ctx, srvReq, srvRsp); err != nil {
		return err
	}
	b, _ := json.Marshal(srvRsp)

	rsp.StatusCode = http.StatusOK
	rsp.Body = string(b)

	return nil
}
