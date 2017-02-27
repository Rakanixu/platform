package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich extracts content from document and add to File
func (gfs *GmailFs) DocEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGmailFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupGmailFile).OptsKazoupFile == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupGmailFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupGmailFile).Modified)
		}

		if f.(*file.KazoupGmailFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = gfs.processDocument(f.(*file.KazoupGmailFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}

// ImgEnrich extracts tags from image and generate thumbnail
func (gfs *GmailFs) ImgEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGmailFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupGmailFile).OptsKazoupFile == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupGmailFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupGmailFile).Modified)
		}

		if f.(*file.KazoupGmailFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = gfs.processImage(gcs, f.(*file.KazoupGmailFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}

// AudioEnrich extracts audio and save it as text
func (gfs *GmailFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGmailFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupGmailFile).OptsKazoupFile == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupGmailFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupGmailFile).Modified)
		}

		if f.(*file.KazoupGmailFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = gfs.processAudio(gcs, f.(*file.KazoupGmailFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}
