package bleve

import (
	"log"
	"sync"

	"github.com/blevesearch/bleve"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
)

type BleveEngine struct {
	IndexPath   string
	Idx         bleve.Index
	mu          sync.Mutex
	Batch       bleve.Batch
	BatchSize   int
	FileChannel chan *crawler.FileMessage
}

func NewBleveEngine(path string) *BleveEngine {
	return &BleveEngine{
		IndexPath:   path,
		Idx:         Init(),
		BatchSize:   10000,
		FileChannel: make(chan *crawler.FileMessage),
	}
}

var (
	IndexPath = "kazoup.idx"
	Batch     = &bleve.Batch{}
)

//Init Bleve serach engine at the start of the service.
func Init() bleve.Index {
	index, err := bleve.Open(IndexPath)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(IndexPath, mapping)
		if err != nil {
			log.Printf("Error creating index : %s", err.Error())
			return nil
		}
		return index
	}
	return index
}

//Batcher
func (be *BleveEngine) Indexer() {
	b := be.Idx.NewBatch()
	go func() {
		for {
			v := <-be.FileChannel
			if b.Size() < be.BatchSize {
				b.Index(v.Id, v.Data)

			} else {
				be.Idx.Batch(b)
				b.Reset()
				s, _ := be.Idx.Stats().MarshalJSON()
				log.Printf("Stats %", string(s))
			}
		}
	}()
}
