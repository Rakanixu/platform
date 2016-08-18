package bleve

import (
	//"encoding/json"
	"errors"
	//"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"log"
	"time"
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
	c := make(chan int)

	go func() {
		for {

			/*			v := <-b.filesChannel
						timeout := <-c

						log.Println("timeout", timeout)
						log.Println("subscriber", v.Id)
						if batch.Size() < b.batchSize {
							batch.Index(v.Id, v.Data)
						} else {
							log.Println("NOw batch")
							b.indexMap[files].Batch(batch)
							batch.Reset()
						}*/
			//c <- batch.Size()

			select {
			case v := <-b.filesChannel:
				log.Print("Ssubscriber", v.Id)
				if batch.Size() < b.batchSize {
					batch.Index(v.Id, v.Data)
				} else {
					log.Println("NOw batch")
					b.indexMap[files].Batch(batch)
					batch.Reset()
				}
			case timeout := <-c:
				if batch.Size() > 0 {
					log.Println("NOw batch", timeout)
					b.indexMap[files].Batch(batch)
					batch.Reset()
				}

				//log.Println("timeout", timeout)
			default:

			}

		}
	}()

	go func() {
		for {
			//sizeBeforeTimeout := <-c
			time.Sleep(time.Second * 2)
			c <- 1
			/*			sizeRealTime := <-c

						log.Println(sizeBeforeTimeout, sizeRealTime)

						if sizeBeforeTimeout == sizeRealTime && sizeRealTime > 0 {
							log.Println("NOw batch")
							b.indexMap[files].Batch(batch)
							batch.Reset()
						} else {
							log.Println("nothing to batch", batch.Size())
						}*/
		}
	}()

	return nil
}
