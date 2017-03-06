package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich extracts content from document and add to File
func (bfs *BoxFs) DocEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupBoxFile)
		if !ok {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching doc file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupBoxFile).OptsKazoupFile == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupBoxFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupBoxFile).Modified)
		}

		if f.(*file.KazoupBoxFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = bfs.processDocument(f.(*file.KazoupBoxFile))
			if err != nil {
				bfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		bfs.FilesChan <- NewFileMsg(f, err)
	}()

	return bfs.FilesChan
}

// ImgEnrich extracts tags from image and generate thumbnail
func (bfs *BoxFs) ImgEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupBoxFile)
		if !ok {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching img file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupBoxFile).OptsKazoupFile == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupBoxFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupBoxFile).Modified)
		}

		if f.(*file.KazoupBoxFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = bfs.processImage(f.(*file.KazoupBoxFile))
			if err != nil {
				bfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		bfs.FilesChan <- NewFileMsg(f, err)
	}()

	return bfs.FilesChan
}

// AudioEnrich extracts audio and save it as text
func (bfs *BoxFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupBoxFile)
		if !ok {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupBoxFile).OptsKazoupFile == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupBoxFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupBoxFile).Modified)
		}

		if f.(*file.KazoupBoxFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = bfs.processAudio(gcs, f.(*file.KazoupBoxFile))
			if err != nil {
				bfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		bfs.FilesChan <- NewFileMsg(f, err)
	}()

	return bfs.FilesChan
}

// ImgEnrich extracts tags from image and generate thumbnail
func (bfs *BoxFs) Thumbnail(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupBoxFile)
		if !ok {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Error generating thumbnail file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processThumb := false
		if f.(*file.KazoupBoxFile).OptsKazoupFile == nil {
			processThumb = true
		} else {
			processThumb = f.(*file.KazoupBoxFile).OptsKazoupFile.ThumbnailTimestamp.Before(f.(*file.KazoupBoxFile).Modified)
		}

		if f.(*file.KazoupBoxFile).Category == globals.CATEGORY_PICTURE && processThumb {
			f, err = bfs.processThumbnail(gcs, f.(*file.KazoupBoxFile))
			if err != nil {
				bfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		bfs.FilesChan <- NewFileMsg(f, err)
	}()

	return bfs.FilesChan
}
