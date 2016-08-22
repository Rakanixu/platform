package elastic

import (
	"log"

	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

var (
	URL = "localhost:9200"
)

type elastic struct {
}

func init() {

	log.Print("Registering ElasticSearch")
	engine.Register(new(elastic))
}

func (e *elastic) Init() error {
	log.Print("Initializing ElasticSearch")
	return nil
}

func (e *elastic) Search(ctx context.Context, req *search.SearchRequest, client client.Client, serviceName string) (*search.SearchResponse, error) {
	return &search.SearchResponse{}, nil
}
