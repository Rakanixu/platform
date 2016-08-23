package bleve

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	lib "github.com/blevesearch/bleve"
	"github.com/kazoup/gabs"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/db/srv/engine"
	db "github.com/kazoup/platform/db/srv/proto/db"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	dataPath        = os.TempDir()
	kazoupNamespace = "/kazoup/"
	filesIndex      = "files"
)

type bleve struct {
	mu           sync.Mutex
	indexMap     map[string]lib.Index
	filesChannel chan *crawler.FileMessage
	batchSize    int
}

func init() {
	engine.Register(&bleve{
		indexMap:     make(map[string]lib.Index),
		filesChannel: make(chan *crawler.FileMessage),
		batchSize:    1000,
	})
}

func (b *bleve) Init() error {
	err := errors.New("")

	files, err := ioutil.ReadDir(os.TempDir() + kazoupNamespace)
	if err != err {
		return err
	}

	b.mu.Lock()
	for _, file := range files {
		if err := openIndex(b, file.Name()); err != nil {
			return err
		}
	}

	// Check files index exists, if not, create. Subscriber needs to have it open
	if b.indexMap[filesIndex] == nil {
		if err := openIndex(b, filesIndex); err != nil {
			return err
		}
	}
	b.mu.Unlock()

	if err := indexer(b); err != nil {
		return err
	}

	return nil
}

func (b *bleve) Subscribe(ctx context.Context, msg *crawler.FileMessage) error {
	b.filesChannel <- msg

	return nil
}

func (b *bleve) Create(req *db.CreateRequest) (*db.CreateResponse, error) {
	response := &db.CreateResponse{}
	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	return response, b.indexMap[req.Index].Index(req.Id, req.Data)
}

func (b *bleve) Read(req *db.ReadRequest) (*db.ReadResponse, error) {
	response := &db.ReadResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	document, err := b.indexMap[req.Index].Document(req.Id)
	if err != nil {
		return response, err
	}

	if document != nil {
		response.Result = string(document.Fields[0].Value())
	}

	return response, err
}

func (b *bleve) Update(req *db.UpdateRequest) (*db.UpdateResponse, error) {
	response := &db.UpdateResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	return response, b.indexMap[req.Index].Index(req.Id, req.Data)
}

func (b *bleve) Delete(req *db.DeleteRequest) (*db.DeleteResponse, error) {
	response := &db.DeleteResponse{}

	if !b.indexExists(req.Index) {
		return response, errors.New("Index does not exists")
	}

	return response, b.indexMap[req.Index].Delete(req.Id)
}

func (b *bleve) CreateIndexWithSettings(req *db.CreateIndexWithSettingsRequest) (*db.CreateIndexWithSettingsResponse, error) {
	err := errors.New("")
	response := &db.CreateIndexWithSettingsResponse{}

	b.mu.Lock()
	if b.indexExists(req.Index) {
		return response, errors.New("Index already exists")
	}

	b.indexMap[req.Index], err = lib.Open(os.TempDir() + kazoupNamespace + req.Index)
	if err != nil {
		mapping := lib.NewIndexMapping()
		b.indexMap[req.Index], err = lib.New(os.TempDir()+kazoupNamespace+req.Index, mapping)
		if err != nil {
			log.Fatalf("Error creating index : %s", err.Error())
			return response, err
		}
	}
	b.mu.Unlock()

	return response, nil
}

func (b *bleve) PutMappingFromJSON(req *db.PutMappingFromJSONRequest) (*db.PutMappingFromJSONResponse, error) {
	response := &db.PutMappingFromJSONResponse{}

	return response, nil
}

func (b *bleve) Status(req *db.StatusRequest) (*db.StatusResponse, error) {
	response := &db.StatusResponse{}

	jsonStatus := gabs.New()
	jsonStatus.SetP("bleve", "master_node")
	jsonStatus.SetP(nil, "metadata.indices")

	for k, _ := range b.indexMap {
		if b.indexExists(k) {
			jsonStatus.SetP("open", "metadata.indices."+k+".state")

		}
	}

	response.Status = jsonStatus.String()

	return response, nil
}

func (b *bleve) Search(req *db.SearchRequest) (*db.SearchResponse, error) {
	var indexSearch, qString string
	var sr *lib.SearchRequest

	if len(req.Index) > 0 {
		indexSearch = req.Index
	} else {
		indexSearch = filesIndex
	}

	if b.indexMap[indexSearch] == nil {
		return &db.SearchResponse{}, errors.New("index does not exists")
	}

	if len(req.Term) > 0 {
		qString += fmt.Sprintf(" %s", req.Term)
	}
	if len(req.Category) > 0 {
		qString += fmt.Sprintf(" +category:%s", req.Category)
	}
	if req.Depth > 0 {
		qString += fmt.Sprintf(" +depth:>=%d +depth:<=%d", req.Depth, req.Depth)
	}
	if len(req.Url) > 0 {
		qString += fmt.Sprintf(" +%s", req.Url)
	}

	log.Println(qString)

	// No fields specify, we want to match all documents in the index
	if len(qString) == 0 {
		q := lib.NewMatchAllQuery()
		sr = &lib.SearchRequest{
			Query: q,
			Size:  int(req.Size),
			From:  int(req.From),
		}
	} else {
		log.Println(qString)
		q := lib.NewQueryStringQuery(qString)
		sr = lib.NewSearchRequestOptions(q, int(req.Size), int(req.From), false)

	}
	sr.Fields = []string{"*"} // Retrieve all fields

	results, err := b.indexMap[indexSearch].Search(sr)
	if err != nil {
		return &db.SearchResponse{}, err
	}

	var buffer bytes.Buffer
	count := 0

	/*	files, _ := json.Marshal(results.Hits)
		log.Println(string(files))*/

	buffer.WriteString(`[`)
	for _, obj := range results.Hits {
		file, _ := json.Marshal(obj.Fields)
		buffer.WriteString(string(file))

		if count < len(results.Hits)-1 {
			buffer.WriteString(`,`)
		}
		count++
	}
	buffer.WriteString(`]`)

	jsonInfo := gabs.New()
	jsonInfo.SetP(results.Total, "total")

	return &db.SearchResponse{
		Result: buffer.String(),
		Info:   jsonInfo.String(),
	}, nil
}
