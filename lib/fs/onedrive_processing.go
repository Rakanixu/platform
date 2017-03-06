package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
)

// DocEnrich extracts content from document and add to File
func (ofs *OneDriveFs) DocEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupOneDriveFile)
		if !ok {
			ofs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Doc)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupOneDriveFile).OptsKazoupFile == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupOneDriveFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupOneDriveFile).Modified)
		}

		if f.(*file.KazoupOneDriveFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = ofs.processDocument(f.(*file.KazoupOneDriveFile))
			if err != nil {
				ofs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		ofs.FilesChan <- NewFileMsg(f, err)
	}()

	return ofs.FilesChan
}

// ImgEnrich extracts tags from image
func (ofs *OneDriveFs) ImgEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupOneDriveFile)
		if !ok {
			ofs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Img)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupOneDriveFile).OptsKazoupFile == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupOneDriveFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupOneDriveFile).Modified)
		}

		if f.(*file.KazoupOneDriveFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = ofs.processImage(f.(*file.KazoupOneDriveFile))
			if err != nil {
				ofs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		ofs.FilesChan <- NewFileMsg(f, err)
	}()

	return ofs.FilesChan
}

// AudioEnrich extracts audio and save it as text
func (ofs *OneDriveFs) AudioEnrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupOneDriveFile)
		if !ok {
			ofs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupOneDriveFile).OptsKazoupFile == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupOneDriveFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupOneDriveFile).Modified)
		}

		if f.(*file.KazoupOneDriveFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = ofs.processAudio(gcs, f.(*file.KazoupOneDriveFile))
			if err != nil {
				ofs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		ofs.FilesChan <- NewFileMsg(f, err)
	}()

	return ofs.FilesChan
}

// Thumbnail generate thumbnail
func (ofs *OneDriveFs) Thumbnail(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupOneDriveFile)
		if !ok {
			ofs.FilesChan <- NewFileMsg(nil, errors.New("Error generating thumbnail file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processThumbnail := false
		if f.(*file.KazoupOneDriveFile).OptsKazoupFile == nil {
			processThumbnail = true
		} else {
			processThumbnail = f.(*file.KazoupOneDriveFile).OptsKazoupFile.ThumbnailTimestamp.Before(f.(*file.KazoupOneDriveFile).Modified)
		}

		if f.(*file.KazoupOneDriveFile).Category == globals.CATEGORY_PICTURE && processThumbnail {
			f, err = ofs.processThumbnail(gcs, f.(*file.KazoupOneDriveFile))
			if err != nil {
				ofs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		ofs.FilesChan <- NewFileMsg(f, err)
	}()

	return ofs.FilesChan
}
