package elastic

import (
	"log"

	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
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

func (e *elastic) Search(req *search.SearchRequest) (res *search.SearchResponse, err error) {
	return nil
}
