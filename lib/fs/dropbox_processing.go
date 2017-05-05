package fs

import (
	"errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
)

// DocEnrich
func (dfs *DropboxFs) DocEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupDropboxFile)
		if !ok {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching doc file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processDoc := false
		if f.(*file.KazoupDropboxFile).OptsKazoupFile == nil || f.(*file.KazoupDropboxFile).OptsKazoupFile.ContentTimestamp == nil {
			processDoc = true
		} else {
			processDoc = f.(*file.KazoupDropboxFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_DOCUMENT && processDoc {
			f, err = dfs.processDocument(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		dfs.FilesChan <- NewFileMsg(f, err)
	}()

	return dfs.FilesChan
}

// ImgEnrich extracts tags from image and generate thumbnail
func (dfs *DropboxFs) ImgEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupDropboxFile)
		if !ok {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching doc file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupDropboxFile).OptsKazoupFile == nil || f.(*file.KazoupDropboxFile).OptsKazoupFile.TagsTimestamp == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupDropboxFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = dfs.processImage(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		dfs.FilesChan <- NewFileMsg(f, err)
	}()

	return dfs.FilesChan
}

// AudioEnrich extracts audio and save it as text
func (dfs *DropboxFs) AudioEnrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupDropboxFile)
		if !ok {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file (Audio)"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.CTagsTimestamp are not defined,
		// Content was never extracted before
		processAudio := false
		if f.(*file.KazoupDropboxFile).OptsKazoupFile == nil || f.(*file.KazoupDropboxFile).OptsKazoupFile.AudioTimestamp == nil {
			processAudio = true
		} else {
			processAudio = f.(*file.KazoupDropboxFile).OptsKazoupFile.AudioTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_AUDIO && processAudio {
			f, err = dfs.processAudio(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		dfs.FilesChan <- NewFileMsg(f, err)
	}()

	return dfs.FilesChan
}

// Thumbnail generates thumbnail
func (dfs *DropboxFs) Thumbnail(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupDropboxFile)
		if !ok {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("Error generating thumbnail file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		processImg := false
		if f.(*file.KazoupDropboxFile).OptsKazoupFile == nil || f.(*file.KazoupDropboxFile).OptsKazoupFile.ThumbnailTimestamp == nil {
			processImg = true
		} else {
			processImg = f.(*file.KazoupDropboxFile).OptsKazoupFile.ThumbnailTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_PICTURE && processImg {
			f, err = dfs.processThumbnail(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		dfs.FilesChan <- NewFileMsg(f, err)
	}()

	return dfs.FilesChan
}
