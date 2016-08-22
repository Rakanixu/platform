package bleve_search

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type bleve struct{}

func init() {
	engine.Register(new(bleve))
}

func (b *bleve) Init() error {
	return nil
}

func (b *bleve) Search(ctx context.Context, req *search.SearchRequest, client client.Client, serviceName string) (*search.SearchResponse, error) {
	srvReq := client.NewRequest(
		serviceName,
		"DB.Search",
		&db_proto.SearchRequest{
			Index:    req.Index,
			Term:     req.Term,
			From:     req.From,
			Size:     req.Size,
			Category: req.Category,
			Url:      req.Url,
			Depth:    req.Depth,
			Type:     req.Type,
		},
	)
	srvRes := &db_proto.SearchResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return &search.SearchResponse{}, err
	}

	return &search.SearchResponse{
		Result: srvRes.Result,
		Info:   srvRes.Info,
	}, nil
}
