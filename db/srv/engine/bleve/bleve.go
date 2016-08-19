package bleve

import (
	"errors"
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
	files           = "files"
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
		b.indexMap[file.Name()], err = lib.Open(os.TempDir() + kazoupNamespace + file.Name())
		if err != nil {
			mapping := lib.NewIndexMapping()
			b.indexMap[file.Name()], err = lib.New(os.TempDir()+kazoupNamespace+file.Name(), mapping)
			if err != nil {
				log.Fatalf("Error creating index : %s", err.Error())
				log.Println("init bleve")
				return err
			}
			return nil
		}
		log.Println("end loop")
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
