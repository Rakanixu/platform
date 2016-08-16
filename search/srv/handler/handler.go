package handler

import (
	"log"

	"github.com/kazoup/platform/search/srv/engine"
	proto "github.com/kazoup/platform/search/srv/proto/search"
	"golang.org/x/net/context"
)

// Search struct
type Search struct {
}

// Create srv handler
func (s *Search) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {

	response, err := engine.Search(req)
	log.Print("handler response : ", response)
	if err != nil {
		return err
	}

	rsp.Result = response.Result
	return nil
	// Instantiate ElasticQuery
	//elasticQueryObj := elastic_query.ElasticQuery{
	//	Term:     req.Term,
	//	From:     req.From,
	//	Size:     req.Size,
	//	Category: req.Category,
	//	Url:      req.Url,
	//	Depth:    req.Depth,
	//	Type:     req.Type,
	//}
	//elasticQuery, err := elasticQueryObj.Query()
	//if err != nil {
	//	return errors.InternalServerError("go.micro.srv.search", err.Error())
	//}

	//// Search in Elastic
	//srvReq := s.Client.NewRequest(
	//	s.SearchServiceName,
	//	"Elastic.Query",
	//	&elastic.QueryRequest{
	//		Index: "files",
	//		Type:  "file",
	//		Query: elasticQuery, // Query for Elastic
	//	},
	//)
	//srvRsp := &elastic.QueryResponse{}

	//if err := s.Client.Call(ctx, srvReq, srvRsp); err != nil {
	//	return errors.InternalServerError("go.micro.srv.search", err.Error())
	//}

	//rsp.Result = srvRsp.Result

	//// TODO: search in other places.. join results, common format, etc..

}
