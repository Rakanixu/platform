package bleve

import (
	//"errors"
	lib "github.com/blevesearch/bleve"
	"github.com/kazoup/platform/db/srv/engine"
	db "github.com/kazoup/platform/db/srv/proto/db"

	"os"
	//"path/filepath"
	"path/filepath"
	"sync"
)

var (
	//IndexFiles = "/tmp/files.idx"
	//IndexDatasources = "/tmp/datasources.idx"
	//IndexFlags       = "/tmp/flags.idx"
	dataPath        = os.TempDir()
	kazoupNamespace = "/kazoup/"
	indexPostfix    = ".idx"
)

type bleve struct {
	idx      lib.Index
	mu       sync.Mutex
	indexMap map[string]lib.Index
}

func init() {
	engine.Register(new(bleve))
}

func (b *bleve) Init() error {

	// Find all indexes folders in dataPath and sstore in dir
	//dir := []byte{}
	filepath.Walk(os.TempDir()+kazoupNamespace, b.walkHandler())

	/*	for _, v := range dir {
		b.idx, err = lib.Open(dataPath + v + indexPostfix)
		if err != nil {

		}
		//append

	}*/

	/*	err := errors.New("")
		b.mu.Lock()

		b.mu.Unlock()
		b.idx, err = lib.Open(IndexFiles)
		if err != nil {
			mapping := lib.NewIndexMapping()
			mapping.StoreDynamic = true

			documentMapping := lib.NewTextFieldMapping()
			documentMapping.Store = true

			mapping.DefaultMapping.AddFieldMappingsAt("data", documentMapping)
			b.idx, err = lib.New(IndexFiles, mapping)
			if err != nil {
				log.Fatalf("Error creating index : %s", err.Error())
				return err
			}
			return nil
		}*/

	return nil
}

func (b *bleve) Create(req *db.CreateRequest) (res *db.CreateResponse, err error) {

	return &db.CreateResponse{}, b.idx.Index(req.Id, req)
}

func (b *bleve) Read(req *db.ReadRequest) (res *db.ReadResponse, err error) {
	response := &db.ReadResponse{}

	document, err := b.idx.Document(req.Id)
	if err != nil {
		return response, err
	}

	if document != nil {
		for _, v := range document.Fields {
			if v.Name() == "data" {
				response.Result = string(v.Value())
			}
		}
	}

	return response, err
}
