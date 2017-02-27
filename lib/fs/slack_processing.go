package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich extracts content from document and add to File
func (sfs *SlackFs) DocEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupSlackFile)
		if !ok {
			sfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Doc)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupSlackFile).OptsKazoupFile == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupSlackFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupSlackFile).Modified)
		}

		if f.(*file.KazoupSlackFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = sfs.processDocument(f.(*file.KazoupSlackFile))
			if err != nil {
				sfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		sfs.FilesChan <- NewFileMsg(f, err)
	}()

	return sfs.FilesChan
}

// ImgEnrich extracts tags from image and generate thumbnail
func (sfs *SlackFs) ImgEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupSlackFile)
		if !ok {
			sfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Img)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupSlackFile).OptsKazoupFile == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupSlackFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupSlackFile).Modified)
		}

		if f.(*file.KazoupSlackFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = sfs.processImage(gcs, f.(*file.KazoupSlackFile))
			if err != nil {
				sfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		sfs.FilesChan <- NewFileMsg(f, err)
	}()

	return sfs.FilesChan
}

// AudioEnrich extracts audio and save it as text
func (sfs *SlackFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupSlackFile)
		if !ok {
			sfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupSlackFile).OptsKazoupFile == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupSlackFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupSlackFile).Modified)
		}

		if f.(*file.KazoupSlackFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = sfs.processAudio(gcs, f.(*file.KazoupSlackFile))
			if err != nil {
				sfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		sfs.FilesChan <- NewFileMsg(f, err)
	}()

	return sfs.FilesChan
}
