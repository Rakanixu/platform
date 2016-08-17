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

		}

		return nil
	}
}
