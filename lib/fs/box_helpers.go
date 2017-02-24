package fs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/kazoup/platform/lib/box"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/image"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"github.com/kennygrant/sanitize"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
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

// processImage, thumbnail generation, cloud vision processing
func (bfs *BoxFs) processImage(gcs *gcslib.GoogleCloudStorage, f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = bcs.Download(f.Original.ID)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		log.Println("ERROR DOWNLOADING FILE", err)
		return nil, err
	}
	defer rc.Close()

	// Split readcloser into two or more for paralel processing
	var buf1, buf2 bytes.Buffer
	w := io.MultiWriter(&buf1, &buf2)

	if _, err = io.Copy(w, rc); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		backoff.Retry(func() error {
			rd, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf2)), globals.THUMBNAIL_WIDTH)
			if err != nil {
				log.Println("THUMNAIL GENERATION ERROR, SKIPPING", err)
				// Skip retry
				return nil
			}

			if err := gcs.Upload(ioutil.NopCloser(rd), bfs.Endpoint.Index, f.ID); err != nil {
				log.Println("THUMNAIL UPLOAD ERROR", err)
				return err
			}

			return nil
		}, backoff.NewExponentialBackOff())
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Resize to optimal size for cloud vision API
		cvrd, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf1)), globals.CLOUD_VISION_IMG_WIDTH)
		if err != nil {
			log.Println("CLOUD VISION ERROR", err)
			return
		}

		if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
			log.Println("CLOUD VISION ERROR", err)
			return
		}

		if f.OptsKazoupFile == nil {
			f.OptsKazoupFile = &file.OptsKazoupFile{
				TagsTimestamp: time.Now(),
			}
		} else {
			f.OptsKazoupFile.TagsTimestamp = time.Now()
		}
	}()

	wg.Wait()

	return f, nil
}

// processDocument sends the original file to tika and enrich KazoupBoxFile with Tika interface
func (bfs *BoxFs) processDocument(f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	rc, err := bcs.Download(f.Original.ID)
	if err != nil {
		return nil, err
	}

	t, err := tika.ExtractContent(rc)
	if err != nil {
		return nil, err
	}

	f.Content = t.Content()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ContentTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.ContentTimestamp = time.Now()
	}

	// Apply rossete
	if len(f.Content) > 0 {
		nl, err := regexp.Compile("\n")
		if err != nil {
			return nil, err
		}
		q, err := regexp.Compile("\"")
		if err != nil {
			return nil, err
		}

		f.Entities, err = rossetelib.Entities(q.ReplaceAllString(nl.ReplaceAllString(sanitize.HTML(f.Content), " "), ""))
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (bfs *BoxFs) processAudio(gcs *gcslib.GoogleCloudStorage, f *file.KazoupBoxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(bfs.Endpoint, globals.Box)
	if err != nil {
		return nil, err
	}

	rc, err := bcs.Download(f.Original.ID)
	if err != nil {
		return nil, err
	}

	if err := gcs.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

	log.Println("111111111111111111")
	stt, err := sttlib.AsyncContent(fmt.Sprintf("gs://%s/%s", globals.AUDIO_BUCKET, f.ID))
	if err != nil {
		return nil, err
	}

	f.Content = stt.Content()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			AudioTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.AudioTimestamp = time.Now()
	}
	log.Println("222222222222222222", f.Content)

	return f, nil
}

// token returns authorization header token as string
func (bfs *BoxFs) token() string {
	return "Bearer " + bfs.Endpoint.Token.AccessToken
}
