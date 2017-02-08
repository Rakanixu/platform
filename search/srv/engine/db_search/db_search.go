package db_search

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
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

func (d *dbSearch) Read(ctx context.Context, req *search.ReadRequest, client client.Client) (*search.ReadResponse, error) {
	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: req.Index,
			Type:  req.Type,
			Id:    req.Id,
		},
	)
	srvRes := &db_proto.ReadResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return &search.ReadResponse{}, err
	}

	return &search.ReadResponse{
		Result: srvRes.Result,
	}, nil
}

func (d *dbSearch) Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error) {
	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Search",
		&db_proto.SearchRequest{
			Index:                req.Index,
			Term:                 req.Term,
			From:                 req.From,
			Size:                 req.Size,
			Category:             req.Category,
			Url:                  req.Url,
			Depth:                req.Depth,
			Type:                 req.Type,
			FileType:             req.FileType,
			Access:               req.Access,
			NoKazoupFileOriginal: req.NoKazoupFileOriginal,
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

func (d *dbSearch) SearchById(ctx context.Context, req *search.SearchByIdRequest, client client.Client) (*search.SearchByIdResponse, error) {
	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.SearchById",
		&db_proto.SearchByIdRequest{
			Index: req.Index,
			Type:  req.Type,
			Id:    req.Id,
		},
	)
	srvRes := &db_proto.SearchByIdResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return &search.SearchByIdResponse{}, err
	}

	return &search.SearchByIdResponse{
		Result: srvRes.Result,
	}, nil
}
