package fs

import (
	"github.com/kazoup/platform/lib/file"
)

func (mfs *MockFs) DocEnrich(f file.File) chan FileMsg {
	// Use a mock file instead of input
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}

func (mfs *MockFs) ImgEnrich(f file.File) chan FileMsg {
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}

func (mfs *MockFs) AudioEnrich(f file.File) chan FileMsg {
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}

func (mfs *MockFs) Thumbnail(f file.File) chan FileMsg {
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}
