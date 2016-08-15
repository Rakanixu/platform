package bleve

import (
	"log"

	"github.com/blevesearch/bleve"
)

type BleveEngine struct {
	IndexPath string
	Idx       bleve.Index
	Batch     bleve.Batch
}

func NewBleveEngine() *BleveEngine {
	return &BleveEngine{
		IndexPath: "kazoup.idx",
		Idx:       Init(),
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
