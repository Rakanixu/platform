package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cenkalti/backoff"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/image"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

// getFiles discover all files in google drive account
func (gfs *GoogleDriveFs) getFiles() error {
	srv, err := gfs.getDriveService()
	if err != nil {
		return err
	}

	r, err := srv.Files.List().PageSize(100).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// getNextPage allows pagination while discovering files
func (gfs *GoogleDriveFs) getNextPage(srv *drive.Service, nextPageToken string) error {
	r, err := srv.Files.List().PageToken(nextPageToken).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// pushFilesToChanForPage sends discovered files to the file system channel
func (gfs *GoogleDriveFs) pushFilesToChanForPage(files []*drive.File) error {
	for _, v := range files {
		f := file.NewKazoupFileFromGoogleDriveFile(*v, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if f != nil {
			gfs.FilesChan <- NewFileMsg(f, nil)
		}
	}

	return nil
}

// processImage, thumbnail generation, cloud vision processing
func (gfs *GoogleDriveFs) processImage(gcs *gcslib.GoogleCloudStorage, f *file.KazoupGoogleFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.GoogleDrive
	gdcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleDrive)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		// Not great, but check implementation for details about variadic params
		rc, err = gdcs.Download(f.Original.Id, "download", "")

		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		log.Println("ERROR DOWNLOADING FILE", err)
		return nil, err
	}
	defer rc.Close()

	// Split readcloser into two or more for different processing
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
			// Resize to our thumbnail size
			rd, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf2)), globals.THUMBNAIL_WIDTH)
			if err != nil {
				log.Println("THUMNAIL GENERATION ERROR, SKIPPING", err)
				// Skip retry
				return nil
			}

			if err := gcs.Upload(ioutil.NopCloser(rd), gfs.Endpoint.Index, f.ID); err != nil {
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

// processDocument sends the original file to tika and enrich KazoupGoogleFile with Tika interface
func (gfs *GoogleDriveFs) processDocument(f *file.KazoupGoogleFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.GoogleDrive
	gcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleDrive)
	if err != nil {
		return nil, err
	}

	// Google documents cannot be dowloaded, they should be exported to required mimeType
	var opts [2]string
	if strings.Contains(f.MimeType, "vnd.google-apps.") {
		opts[0] = "export"
		opts[1] = globals.GoogleDriveExportAs(f.MimeType)
	} else {
		opts[0] = "download"
		opts[1] = ""
	}

	rc, err := gcs.Download(f.Original.Id, opts[0], opts[1])
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

	return f, nil
}

// processAudio uploads audio file to GCS and runs async speech to text over it
func (gfs *GoogleDriveFs) processAudio(gcs *gcslib.GoogleCloudStorage, f *file.KazoupGoogleFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.GoogleDrive
	gdcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleDrive)
	if err != nil {
		return nil, err
	}

	// Google documents cannot be dowloaded, they should be exported to required mimeType
	var opts [2]string
	if strings.Contains(f.MimeType, "vnd.google-apps.") {
		opts[0] = "export"
		opts[1] = globals.GoogleDriveExportAs(f.MimeType)
	} else {
		opts[0] = "download"
		opts[1] = ""
	}

	rc, err := gdcs.Download(f.Original.Id, opts[0], opts[1])
	if err != nil {
		return nil, err
	}

	if err := gcs.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

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

	return f, nil
}

// getDriveService return a google drive service instance
func (gfs *GoogleDriveFs) getDriveService() (*drive.Service, error) {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	return drive.New(c)
}
