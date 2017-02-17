package fs

import (
	"bufio"
	"bytes"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/tika"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"io"
	"io/ioutil"
	"log"
	"strings"
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
			if err := gfs.processImage(f, f.ID); err != nil {
				log.Println(err)
			}

			if err := gfs.processDocument(f); err != nil {
				log.Println(err)
			}

			gfs.FilesChan <- NewFileMsg(f, nil)
		}
	}

	return nil
}

func (gfs *GoogleDriveFs) processImage(f *file.KazoupGoogleFile, id string) error {
	c := categories.GetDocType("." + f.Original.FullFileExtension)
	if len(f.Original.FullFileExtension) == 0 {
		c = categories.GetDocType(f.MimeType)
	}

	if c == globals.CATEGORY_PICTURE {
		// Download file from GoogleDrive, so connector is globals.GoogleDrive
		gcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleDrive)
		if err != nil {
			return err
		}

		// Not great, but check implementation for details about variadic params
		rc, err := gcs.Download(f.Original.Id, "download", "")
		if err != nil {
			return err
		}

		// Split readcloser into two or more for different processing
		var buf1, buf2 bytes.Buffer
		w := io.MultiWriter(&buf1, &buf2)

		if _, err = io.Copy(w, rc); err != nil {
			return err
		}

		go func() {
			// Resize to our thumbnail size
			rd, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf2)), globals.THUMBNAIL_WIDTH)
			if err != nil {
				log.Println(err)
			}

			// Upload file to GoogleCloudStorage, so connector is globals.GoogleCloudStorage
			ncs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleCloudStorage)
			if err != nil {
				log.Println(err)
			}

			if err := ncs.Upload(rd, id); err != nil {
				log.Println(err)
			}
		}()

		// Resize to optimal size for cloud vision API
		cvrd, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf1)), globals.CLOUD_VISION_IMG_WIDTH)
		if err != nil {
			return err
		}

		if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
			return err
		}

		f.OptsKazoupFile.TagsTimestamp = time.Now()
	}

	return nil
}

// processDocument sends the original file to tika and enrich KazoupGoogleFile with Tika interface
func (gfs *GoogleDriveFs) processDocument(f *file.KazoupGoogleFile) error {
	if f.Category == globals.CATEGORY_DOCUMENT {
		// Download file from GoogleDrive, so connector is globals.GoogleDrive
		gcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.GoogleDrive)
		if err != nil {
			return err
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
			return err
		}

		t, err := tika.ExtractContent(rc)
		if err != nil {
			return err
		}

		f.Content = t.Content()
	}

	return nil
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
