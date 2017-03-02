package fs

import (
	"github.com/kazoup/platform/lib/file"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich
func (lfs *LocalFs) DocEnrich(f file.File) chan FileMsg {
	return lfs.FilesChan
}

// ImgEnrich
func (lfs *LocalFs) ImgEnrich(f file.File) chan FileMsg {
	return lfs.FilesChan
}

// AudioEnrich
func (lfs *LocalFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	return lfs.FilesChan
}

// ImgEnrich
func (lfs *LocalFs) Thumbnail(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	return lfs.FilesChan
}
