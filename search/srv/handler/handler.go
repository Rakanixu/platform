package handler

import (
	"github.com/kazoup/platform/search/srv/engine"
	proto "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	//"github.com/micro/go-micro/errors"
	//"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	//"log"
)

// Search struct
type Search struct {
	Client client.Client
}

// Search srv handler
func (s *Search) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	/*	md, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.InternalServerError("search.search", "Unable to retrieve metadata")
		}
		log.Println("CONTEXT")
		log.Println(md["User"])
		log.Println(md["Token"])*/

	response, err := engine.Search(ctx, req, s.Client)
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
