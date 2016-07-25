package handler

import (
	"encoding/json"
	"net/http"

	datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	api "github.com/micro/micro/api/proto"
	"golang.org/x/net/context"
)

// Datasource struct
type Datasource struct{}

// Create datasource handler
func (ds *Datasource) Create(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	// Required fields
	if input["endpoint"] == nil {
		return errors.BadRequest("go.micro.api.datasource", "endpoint required")
	}

	if input["endpoint"].(map[string]interface{})["url"] == nil {
		return errors.BadRequest("go.micro.api.datasource", "endpoint url required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.datasource",
		"DataSource.Create",
		&datasource.CreateRequest{
			Endpoint: &datasource.Endpoint{
				Url: input["endpoint"].(map[string]interface{})["url"].(string),
			},
		},
	)
	srvRes := &datasource.CreateResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	rsp.Body = `{}`
	rsp.StatusCode = http.StatusOK

	return nil
}

// Delete datasource handler
func (ds *Datasource) Delete(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	// Required fields
	if input["id"] == nil {
		return errors.BadRequest("go.micro.api.datasource", "id required")
	}

	srvReq := client.NewRequest(
		"go.micro.srv.datasource",
		"DataSource.Delete",
		&datasource.DeleteRequest{
			Id: input["id"].(string),
		},
	)
	srvRes := &datasource.DeleteResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	rsp.Body = `{}`
	rsp.StatusCode = http.StatusOK

	return nil
}

// Search datasource handler
func (ds *Datasource) Search(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var input map[string]interface{}
	var offset, limit int64
	var query string

	if err := json.Unmarshal([]byte(req.Body), &input); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	// Get args if exists
	if input["query"] != nil {
		query = input["query"].(string)
	}

	if input["offset"] != nil {
		offset = int64(input["offset"].(float64))
	}

	if input["limit"] != nil {
		limit = int64(input["limit"].(float64))
	}

	srvReq := client.NewRequest(
		"go.micro.srv.datasource",
		"DataSource.Search",
		&datasource.SearchRequest{
			Query:  query,
			Offset: offset,
			Limit:  limit,
		},
	)
	srvRes := &datasource.SearchResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	b, err := json.Marshal(srvRes.Result)
	if err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	rsp.Body = string(b)
	rsp.StatusCode = http.StatusOK

	return nil
}

// Scan datasource handler
func (ds *Datasource) Scan(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var scanRequest *datasource.ScanRequest

	if err := json.Unmarshal([]byte(req.Body), &scanRequest); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	srvReq := client.NewRequest(
		"go.micro.srv.datasource",
		"DataSource.Scan",
		scanRequest,
	)
	srvRes := &datasource.ScanResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.InternalServerError("go.micro.api.datasource", err.Error())
	}

	rsp.Body = `{}`
	rsp.StatusCode = http.StatusOK

	return nil
}
