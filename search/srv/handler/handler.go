package handler

import (
	"github.com/kazoup/platform/search/srv/engine"
	proto "github.com/kazoup/platform/search/srv/proto/search"
	"golang.org/x/net/context"
)

// Search struct
type Search struct{}

// Create srv handler
func (s *Search) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {

	response, err := engine.Search(req)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	rsp.Info = response.Info
	return nil
}
