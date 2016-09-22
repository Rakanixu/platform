package db_search

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type dbSearch struct{}

func init() {
	engine.Register(new(dbSearch))
}

func (d *dbSearch) Init() error {
	return nil
}

func (d *dbSearch) Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error) {
	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
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
			FileType: req.FileType,
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

func (d *dbSearch) Aggregate(ctx context.Context, req *search.AggregateRequest, client client.Client) (*search.AggregateResponse, error) {
	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Aggregate",
		req,
	)
	srvRes := &search.AggregateResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return nil, err
	}

	return srvRes, nil
}
