package local

import "os"

// LocalFile model
type LocalFile struct {
	Path string
	Info os.FileInfo
}
