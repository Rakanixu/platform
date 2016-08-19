package bleve

import (
	"errors"
	"fmt"
	lib "github.com/blevesearch/bleve"
	"github.com/kazoup/platform/search/srv/engine"
	search "github.com/kazoup/platform/search/srv/proto/search"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	dataPath        = os.TempDir()
	kazoupNamespace = "/kazoup/"
	files           = "files"
)

type bleve struct {
	mu       sync.Mutex
	indexMap map[string]lib.Index
}

func init() {
	engine.Register(new(bleve))
}

func (b *bleve) Init() error {
	err := errors.New("")

	files, _ := ioutil.ReadDir(os.TempDir() + kazoupNamespace)

	b.mu.Lock()
	for _, file := range files {
		b.indexMap[file.Name()], err = lib.Open(os.TempDir() + kazoupNamespace + file.Name())
		if err != nil {
			mapping := lib.NewIndexMapping()
			b.indexMap[file.Name()], err = lib.New(os.TempDir()+kazoupNamespace+file.Name(), mapping)
			if err != nil {
				log.Fatalf("Error creating index : %s", err.Error())
				return err
			}
			return nil
		}
	}
	b.mu.Unlock()

	return nil
}

func (b *bleve) Search(req *search.SearchRequest) (*search.SearchResponse, error) {
	var indexSearch string

	if len(req.Index) > 0 {
		indexSearch = req.Index
	} else {
		indexSearch = files
	}

	if b.indexMap[indexSearch] == nil {
		return &search.SearchResponse{}, errors.New("index does not exists")
	}

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

	q := lib.NewQueryStringQuery(qString)
	br := lib.NewSearchRequest(q)
	br.Highlight = lib.NewHighlightWithStyle("html")

	results, err := b.indexMap[indexSearch].Search(br)
	if err != nil {
		log.Print(err.Error())
		return &search.SearchResponse{}, nil
	}
	log.Print(results.String())

	return &search.SearchResponse{
		Result: results.String(),
		Info:   "",
	}, nil
}
