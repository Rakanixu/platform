package fs

import (
	"github.com/kazoup/platform/lib/file"
)

type FileMeta struct {
	File  file.File
	Error error
}

func NewFileMeta(file file.File, err error) FileMeta {
	return FileMeta{
		File:  file,
		Error: err,
	}
}
