package handler

import (
	"github.com/kazoup/platform/search/srv/engine"
	proto "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Search struct
type Search struct {
	Client             client.Client
	ElasticServiceName string
}

// Create srv handler
func (s *Search) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	response, err := engine.Search(ctx, req, s.Client, s.ElasticServiceName)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	rsp.Info = response.Info
	return nil
}
