package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich extracts content from document and add to File
func (gfs *GoogleDriveFs) DocEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGoogleFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Doc)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupGoogleFile).OptsKazoupFile == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupGoogleFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
		}

		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = gfs.processDocument(f.(*file.KazoupGoogleFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}

// ImgEnrich extracts tags from image
func (gfs *GoogleDriveFs) ImgEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGoogleFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Img)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupGoogleFile).OptsKazoupFile == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupGoogleFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
		}

		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = gfs.processImage(f.(*file.KazoupGoogleFile))
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
func (gfs *GoogleDriveFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGoogleFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupGoogleFile).OptsKazoupFile == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupGoogleFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
		}

		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = gfs.processAudio(gcs, f.(*file.KazoupGoogleFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}

// Thumbnail generate thumbnail
func (gfs *GoogleDriveFs) Thumbnail(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGoogleFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Img)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processThumbnail := false
		if f.(*file.KazoupGoogleFile).OptsKazoupFile == nil {
			processThumbnail = true
		} else {
			processThumbnail = f.(*file.KazoupGoogleFile).OptsKazoupFile.ThumbnailTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
		}

		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_PICTURE && processThumbnail {
			f, err = gfs.processThumbnail(gcs, f.(*file.KazoupGoogleFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}
