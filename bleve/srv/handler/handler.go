package handler

import (
	"github.com/blevesearch/bleve"
	proto "github.com/kazoup/platform/bleve/srv/proto/bleve"

	"golang.org/x/net/context"
)

type SearchBleve struct {
	Index bleve.Index
	Batch bleve.Batch
}

func (bs *SearchBleve) Search(ctx context.Context, req *proto.SearchRequest, res *proto.SearchResponse) error {

	query := bleve.NewMatchQuery(req.Term)
	request := bleve.NewSearchRequest(query)
	results, err := bs.Index.Search(request)
	if err != nil {
		return err
	}
	res.Result = results.String()
	return nil
}

func (bs *SearchBleve) Create(ctx context.Context, req *proto.CreateRequest, res *proto.CreateResponse) error {
	if err := bs.Index.Index(req.Id, req.Data); err != nil {
		return err
	}
	res.Status = "OK"
	return nil
}
