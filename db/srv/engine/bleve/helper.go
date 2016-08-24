package bleve

import (
	"encoding/json"
	"errors"
	lib "github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzers/keyword_analyzer"
	_ "github.com/blevesearch/bleve/index/store/goleveldb"
	"github.com/kazoup/platform/structs"
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
					file := new(structs.DesktopFile)

					if err := json.Unmarshal([]byte(v.Data), file); err != nil {
						log.Printf("%v", err)
					}

					batch.Index(v.Id, file)
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
		mapping.DefaultAnalyzer = keyword_analyzer.Name
		kvconfig := make(map[string]interface{})

		b.indexMap[indexName], err = lib.NewUsing(os.TempDir()+kazoupNamespace+indexName, mapping, "upside_down", "goleveldb", kvconfig)
		if err != nil {
			return err
		}
	}

	return nil
}
