package bleve

import (
	"errors"
	"fmt"
	"log"

	lib "github.com/blevesearch/bleve"
	"github.com/kazoup/platform/search/srv/engine"
	"github.com/kazoup/platform/search/srv/proto/search"
)

var (
	IndexPath = "/tmp/files.idx"
)

type bleve struct {
	idx lib.Index
}

func init() {
	log.Print("Registering bleve")
	engine.Register(new(bleve))
}

func (b *bleve) Init() error {
	log.Print("Initializing bleve")
	err := errors.New("")
	b.idx, err = lib.Open(IndexPath)
	if err != nil {

		log.Print("Index doesnt exists creating new one")
		mapping := lib.NewIndexMapping()
		b.idx, err = lib.New(IndexPath, mapping)
		if err != nil {
			log.Printf("Error creating index : %s", err.Error())
			return err
		}
		return nil
	}

	log.Print("Succefuly open bleve index")
	return nil
}

func (b *bleve) Search(req *search.SearchRequest) (res *search.SearchResponse, err error) {
	log.Print("New querr %s", req.Term)
	qString := ""
	if len(req.Term) > 0 {
		qString += fmt.Sprintf("%s", req.Term)
	}
	if len(req.Category) > 0 {
		qString += fmt.Sprintf("+category:%s", req.Category)
	}
	if req.Depth != 0 {
		qString += fmt.Sprintf("+depth:%s", string(req.Depth))
	}

	log.Printf("Query %s", qString)
	q := lib.NewQueryStringQuery(qString)
	br := lib.NewSearchRequest(q)
	br.Highlight = lib.NewHighlightWithStyle("html")
	//br.Fields = []string{"name","description"}
	results, err := b.idx.Search(br)
	if err != nil {
		log.Print(err.Error())
		return &search.SearchResponse{}, nil
	}
	log.Print(results.String())
	res = &search.SearchResponse{
		Result: results.String(),
	}
	return res, nil
}
