package handler

import (
	"encoding/json"
	"fmt"
	db "github.com/kazoup/platform/db/srv/proto/db"
	proto "github.com/kazoup/platform/flag/srv/proto/flag"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"strconv"
)

// Flag struct
type Flag struct {
	Client        client.Client
	DbServiceName string
}

// Create srv handler
func (f *Flag) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	if len(req.Key) <= 0 || len(req.Description) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Create", "Fields required")
	}

	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	srvReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Create",
		&db.CreateRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Id for flags will be the key given, so we can RUD easily
			Data:  string(data),
		},
	)
	srvRsp := &db.CreateResponse{}
	if err := f.Client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Create", err.Error())
	}

	return nil
}

// Read srv handler
func (f *Flag) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	if len(req.Key) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Read", "Flag key required")
	}

	srvReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Read",
		&db.ReadRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvRsp := &db.ReadResponse{}
	if err := f.Client.Call(ctx, srvReq, srvRsp); err != nil {
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
	srvReadReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Read",
		&db.ReadRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvReadRsp := &db.ReadResponse{}
	if err := f.Client.Call(ctx, srvReadReq, srvReadRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Flip", err.Error())
	}

	// Flip value
	var input map[string]interface{}

	json.Unmarshal([]byte(srvReadRsp.Result), &input)
	value, _ := strconv.ParseBool(fmt.Sprintf("%v", input["value"]))
	input["value"] = !value
	data, _ := json.Marshal(input)

	// Update the record
	srvUpdateReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Update",
		&db.UpdateRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
			Data:  string(data),
		},
	)
	srvUpdateRsp := &db.UpdateResponse{}
	if err := f.Client.Call(ctx, srvUpdateReq, srvUpdateRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Flip", err.Error())
	}

	return nil
}

// Delete srv handler
func (f *Flag) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Key) <= 0 {
		return errors.BadRequest("go.micro.srv.flag.Flag.Delete", "Flag key required")
	}

	srvReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Delete",
		&db.DeleteRequest{
			Index: "flags", // Hardcoded index for flags
			Type:  "flag",  // Hardcoded type ...
			Id:    req.Key, // Our ID for flags index
		},
	)
	srvRsp := &db.DeleteResponse{}
	if err := f.Client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.Delete", err.Error())
	}

	return nil
}

// List srv handler
func (f *Flag) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	// Search in flags index
	srvReq := f.Client.NewRequest(
		f.DbServiceName,
		"DB.Search",
		&db.SearchRequest{
			Index: "flags",
			Type:  "flag",
			From:  0,
			Size:  9999,
		},
	)
	srvRsp := &db.SearchResponse{}
	if err := f.Client.Call(ctx, srvReq, srvRsp); err != nil {
		return errors.InternalServerError("go.micro.srv.flag.Flag.List", err.Error())
	}

	rsp.Result = srvRsp.Result

	return nil
}
