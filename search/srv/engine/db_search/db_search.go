package db_search

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
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
	srvRes, err := db_helper.ReadFromDB(client, ctx, &db_proto.ReadRequest{
		Index: req.Index,
		Type:  req.Type,
		Id:    req.Id,
	})
	if err != nil {
		return &search.ReadResponse{}, err
	}

	return &search.ReadResponse{
		Result: srvRes.Result,
	}, nil
}

func (d *dbSearch) Search(ctx context.Context, req *search.SearchRequest, client client.Client) (*search.SearchResponse, error) {
	srvRes, err := db_helper.SearchFromDB(client, ctx, &db_proto.SearchRequest{
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
		ContentCategory:      req.ContentCategory,
		NoKazoupFileOriginal: req.NoKazoupFileOriginal,
	})
	if err != nil {
		return &search.SearchResponse{}, err
	}

	return &search.SearchResponse{
		Result: srvRes.Result,
		Info:   srvRes.Info,
	}, nil
}

func (d *dbSearch) SearchById(ctx context.Context, req *search.SearchByIdRequest, client client.Client) (*search.SearchByIdResponse, error) {
	srvRes, err := db_helper.SearchById(client, ctx, &db_proto.SearchByIdRequest{
		Index: req.Index,
		Type:  req.Type,
		Id:    req.Id,
	})
	if err != nil {
		return &search.SearchByIdResponse{}, err
	}

	return &search.SearchByIdResponse{
		Result: srvRes.Result,
	}, nil
}
