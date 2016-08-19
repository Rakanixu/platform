package bleve

import (
	//"encoding/json"
	"errors"
	//"github.com/kazoup/platform/crawler/srv/proto/crawler"
	lib "github.com/blevesearch/bleve"
	"log"
	//"time"
)

func (b *bleve) indexExists(index string) bool {
	if b.indexMap[index] == nil {
		return false
	}
	return true
}

func indexer(b *bleve) error {
	if !b.indexExists(files) {
		return errors.New("index does not exist")
	}

	batch := b.indexMap[files].NewBatch()
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
	b.indexMap[files].Batch(batch)
	batch.Reset()
}
