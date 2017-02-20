package fs

import (
	"bufio"
	"bytes"
	"encoding/json"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/onedrive"
	"github.com/kazoup/platform/lib/tika"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	FACETS = "shared,id,name,size,parentReference,createdBy,fileSystemInfo,lastModifiedDateTime,lastModifiedBy,webUrl,file,folder"
)

// getFiles retrieves drives, directories and files
func (ofs *OneDriveFs) getFiles() error {
	if err := ofs.getDrives(); err != nil {
		return err
	}
	if err := ofs.getDrivesChildren(); err != nil {
		return err
	}

	return nil
}

// getDrives retrieve user drives
func (ofs *OneDriveFs) getDrives() error {
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := globals.OneDriveEndpoint + Drives
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.token())
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var drivesRsp *onedrive.DrivesListResponse
	if err := json.NewDecoder(res.Body).Decode(&drivesRsp); err != nil {
		return err
	}

	for _, v := range drivesRsp.Value {
		ofs.DrivesId = append(ofs.DrivesId, v.ID)
	}

	return nil
}

// getDrivesChildren gets first level element from every found  drive
func (ofs *OneDriveFs) getDrivesChildren() error {
	var url string
	c := &http.Client{}

	for _, v := range ofs.DrivesId {
		//https://api.onedrive.com/v1.0/drives/f5a34c5d0f17415a/root/children
		url = globals.OneDriveEndpoint + Drives + v + "/root/children?select=" + FACETS

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", ofs.token())
		if err != nil {
			return err
		}
		res, err := c.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var filesRsp *onedrive.FilesListResponse
		if err := json.NewDecoder(res.Body).Decode(&filesRsp); err != nil {
			return err
		}

		for _, v := range filesRsp.Value {
			// Is directory
			if len(v.File.MimeType) == 0 {
				ofs.Directories <- v.ID
				// Is file
			} else {
				if err := ofs.pushToFilesChannel(v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getDirChildren get children from directory
func (ofs *OneDriveFs) getDirChildren(id string) error {
	// https://api.onedrive.com/v1.0/drive/items/F5A34C5D0F17415A!114/children
	c := &http.Client{}
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/children?select=" + FACETS
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.token())
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var filesRsp *onedrive.FilesListResponse
	if err := json.NewDecoder(res.Body).Decode(&filesRsp); err != nil {
		return err
	}

	for _, v := range filesRsp.Value {
		if len(v.File.MimeType) == 0 {
			ofs.Directories <- v.ID
		} else {
			if err := ofs.pushToFilesChannel(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ofs *OneDriveFs) getPermisions(f *file.KazoupOneDriveFile) error {
	c := &http.Client{}
	url := globals.OneDriveEndpoint + Drive + "items/" + f.Original.ID + "/permissions"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.token())
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var pRsp *onedrive.PermissionsResponse
	if err := json.NewDecoder(res.Body).Decode(&pRsp); err != nil {
		return err
	}

	for _, v := range pRsp.Value {
		if v.GrantedTo == nil {
			f.Original.PublicURL = v.Link.WebURL
			f.Access = globals.ACCESS_PUBLIC
			break
		}
	}

	return nil
}

// pushToFilesChannel
func (ofs *OneDriveFs) pushToFilesChannel(f onedrive.OneDriveFile) error {
	kof := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)

	if err := ofs.getPermisions(kof); err != nil {
		log.Println(err)
	}

	ofs.FilesChan <- NewFileMsg(kof, nil)

	return nil
}

// processImage, thumbnail generation, cloud vision processing
func (ofs *OneDriveFs) processImage(gcs *gcslib.GoogleCloudStorage, f *file.KazoupOneDriveFile) (file.File, error) {
	// Download file from OneDrive, so connector is globals.OneDrive
	ocs, err := cs.NewCloudStorageFromEndpoint(ofs.Endpoint, globals.OneDrive)
	if err != nil {
		return nil, err
	}

	rc, err := ocs.Download(f.Original.ID)
	if err != nil {
		return nil, err
	}

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

		b, err := image.Thumbnail(ioutil.NopCloser(bufio.NewReader(&buf2)), globals.THUMBNAIL_WIDTH)
		if err != nil {
			log.Println(err)
			return
		}

		if err := gcs.Upload(b, ofs.Endpoint.Index, f.ID); err != nil {
			log.Println(err)
			return
		}
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

// processDocument sends the original file to tika and enrich KazoupOneDriveFile with Tika interface
func (ofs *OneDriveFs) processDocument(f *file.KazoupOneDriveFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.OneDrive
	ocs, err := cs.NewCloudStorageFromEndpoint(ofs.Endpoint, globals.OneDrive)
	if err != nil {
		return nil, err
	}

	rc, err := ocs.Download(f.Original.ID)
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

// token returns authorization header
func (ofs *OneDriveFs) token() string {
	return ofs.Endpoint.Token.TokenType + " " + ofs.Endpoint.Token.AccessToken
}
