package bleve

import (
	"bytes"
	"encoding/gob"
	"errors"
	lib "github.com/blevesearch/bleve"
	"log"
	"os"
)

func (b *bleve) indexExists(index string) bool {
	if b.indexMap[index] == nil {
		return false
	}
	return true
}

func indexer(b *bleve) error {
	if !b.indexExists(filesIndex) {
		return errors.New("index does not exist")
	}

	batch := b.indexMap[filesIndex].NewBatch()
	//c := make(chan int)

	go func() {
		for {
			select {
			case v := <-b.filesChannel:
				log.Print("Ssubscriber", v.Id)
				if batch.Size() < b.batchSize {
					batch.Index(v.Id, v.Data)
				} else {
					log.Println("NOw batch")
					doBatch(b, batch)
				}
				/*		case timeout := <-c:
						if batch.Size() > 0 {
							log.Println("NOw batch", timeout)
							doBatch(b, batch)
						}*/

			}

		}
	}()

	/*	go func() {
		for {
			time.Sleep(time.Second * 10)
			c <- 1
		}
	}()*/

	return nil
}

func doBatch(b *bleve, batch *lib.Batch) {
	b.indexMap[filesIndex].Batch(batch)
	batch.Reset()
}

func openIndex(b *bleve, indexName string) error {
	err := errors.New("")
	b.indexMap[indexName], err = lib.Open(os.TempDir() + kazoupNamespace + indexName)
	if err != nil {
		mapping := lib.NewIndexMapping()
		b.indexMap[indexName], err = lib.New(os.TempDir()+kazoupNamespace+indexName, mapping)
		if err != nil {
			return err
		}
	}

	return nil
}
