package fs

import (
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/kazoup/platform/lib/box"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/objectstorage"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// getDirChildren get children from directory
func (bfs *BoxFs) getDirChildren(id string, offset, limit int) error {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/?offset=%d&limit=%d", globals.BoxFoldersEndpoint, id, offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var bdc *box.BoxDirContents
	if err := json.NewDecoder(rsp.Body).Decode(&bdc); err != nil {
		return err
	}

	for _, v := range bdc.ItemCollection.Entries {
		if v.Type == "folder" {
			// Push found directories into the queue to be crawled
			bfs.Directories <- v.ID
		} else {
			// File discovered, but need to retrieve more info about the file
			if err := bfs.getMetadataFromFile(v.ID); err != nil {
				return err
			}
		}
	}

	if bdc.ItemCollection.TotalCount > bdc.ItemCollection.Offset+bdc.ItemCollection.Limit {
		bfs.getDirChildren(
			id,
			bdc.ItemCollection.Offset+bdc.ItemCollection.Limit,
			bdc.ItemCollection.Limit,
		)
	}

	return nil
}

// getMetadataFromFile retrieves more info about discovered files in box
func (bfs *BoxFs) getMetadataFromFile(id string) error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxFileMetadataEndpoint+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var fm *box.BoxFileMeta
	if err := json.NewDecoder(rsp.Body).Decode(&fm); err != nil {
		return err
	}

	f := file.NewKazoupFileFromBoxFile(*fm, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

	bfs.FilesChan <- NewFileMsg(f, nil)

	return nil
}

// processImage,  cloud vision processing
func (bfs *BoxFs) processImage(f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = bcs.Download(f.OriginalID)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	// Resize to optimal size for cloud vision API
	cvrd, err := image.Thumbnail(rc, globals.CLOUD_VISION_IMG_WIDTH)
	if err != nil {
		return nil, err
	}

	if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
		return nil, err
	}

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			TagsTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.TagsTimestamp = &n
	}

	return f, nil
}

// processThumbnail, thumbnail generation
func (bfs *BoxFs) processThumbnail(f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = bcs.Download(f.OriginalID)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	backoff.Retry(func() error {
		rd, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			// Skip retry
			return nil
		}

		if err := objectstorage.Upload(ioutil.NopCloser(rd), bfs.Endpoint.Index, f.ID); err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ThumbnailTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ThumbnailTimestamp = &n
	}

	return f, nil
}

// processDocument sends the original file to tika and enrich KazoupBoxFile with Tika interface
func (bfs *BoxFs) processDocument(f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	rc, err := bcs.Download(f.OriginalID)
	if err != nil {
		return nil, err
	}

	t, err := tika.ExtractPlainContent(rc)
	if err != nil {
		return nil, err
	}

	f.Content = t.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ContentTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ContentTimestamp = &n
	}

	return f, nil
}

// processAudio uploads audio file to GCS and runs async speech to text over it
func (bfs *BoxFs) processAudio(f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	rc, err := bcs.Download(f.OriginalID)
	if err != nil {
		return nil, err
	}

	if err := objectstorage.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

	stt, err := sttlib.AsyncContent(fmt.Sprintf("gs://%s/%s", globals.AUDIO_BUCKET, f.ID))
	if err != nil {
		return nil, err
	}

	f.Content = stt.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			AudioTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.AudioTimestamp = &n
	}

	return f, nil
}

// token returns authorization header token as string
func (bfs *BoxFs) token() string {
	return "Bearer " + bfs.Endpoint.Token.AccessToken
}
