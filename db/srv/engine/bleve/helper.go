package bleve

import (
	"log"
	"os"
	"path/filepath"
)

func (b *bleve) walkHandler() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		select {
		default:
			log.Println(path)

			err := errors.New("")
			b.mu.Lock()

			
			b.indexMap[path], err = lib.Open(path)
			if err != nil {
		/*	mapping := lib.NewIndexMapping()
			mapping.StoreDynamic = true

			documentMapping := lib.NewTextFieldMapping()
			documentMapping.Store = true

			mapping.DefaultMapping.AddFieldMappingsAt("data", documentMapping)*/
			b.idx, err = lib.New(path, lib.NewIndexMapping())
			if err != nil {
				log.Fatalf("Error creating index : %s", err.Error())
				return err
			}

			b.mu.Unlock()
			return nil
		}


		}

		return nil
	}
}
