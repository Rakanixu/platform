package handler

import (
	"github.com/kazoup/platform/search/srv/engine"
	proto "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Search struct
type Search struct {
	Client client.Client
}

// Read srv handler - proxy to DB.Read
func (s *Search) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	response, err := engine.Read(ctx, req, s.Client)
	if err != nil {
		return err
	}

	rsp.Result = response.Result

	return nil
}

// Search srv handler
func (s *Search) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	response, err := engine.Search(ctx, req, s.Client)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}

// SearchById srv handler
func (s *Search) SearchById(ctx context.Context, req *proto.SearchByIdRequest, rsp *proto.SearchByIdResponse) error {
	response, err := engine.SearchById(ctx, req, s.Client)
	if err != nil {
		return err
	}

	rsp.Result = response.Result

	return nil
}

func (s *Search) Health(ctx context.Context, req *proto.HealthRequest, rsp *proto.HealthResponse) error {
	rsp.Status = 200
	rsp.Info = "OK"

	return nil
}
