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

// Search srv handler
func (s *Search) SearchProxy(ctx context.Context, req *proto.SearchProxyRequest, rsp *proto.SearchProxyResponse) error {
	response, err := engine.SearchProxy(ctx, req, s.Client)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}

// Aggregate srv handler
func (s *Search) Aggregate(ctx context.Context, req *proto.AggregateRequest, rsp *proto.AggregateResponse) error {
	response, err := engine.Aggregate(ctx, req, s.Client)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	rsp.Info = response.Info

	return nil
}
