package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/onedrive"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	Drive  = "drive/"
	Drives = "drives/"
)

// OneDriveFs one drive file system
type OneDriveFs struct {
	Endpoint     *datasource_proto.Endpoint
	Running      chan bool
	FilesChan    chan file.File
	FileMetaChan chan FileMeta
	DrivesId     []string
	Directories  chan string
	LastDirTime  int64
}

// NewOneDriveFsFromEndpoint constructor
func NewOneDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &OneDriveFs{
		Endpoint:     e,
		Running:      make(chan bool, 1),
		FilesChan:    make(chan file.File),
		FileMetaChan: make(chan FileMeta),
		DrivesId:     []string{},
		// This is important to have a size bigger than one, the bigger, less likely to block
		// If not, program execution will block, due to recursivity,
		// We are pushing more elements before finish execution.
		// I expect to never push 10000 folders before other folders have been completly scanned
		Directories: make(chan string, 10000),
	}
}

// List returns 2 channels, for files and state. Discover files in one drive datasources
func (ofs *OneDriveFs) List(c client.Client) (chan file.File, chan bool, error) {
	go func() {
		ofs.LastDirTime = time.Now().Unix()
		for {
			select {
			case v := <-ofs.Directories:
				ofs.LastDirTime = time.Now().Unix()

				err := ofs.getDirChildren(c, v)
				if err != nil {
					log.Println(err)
				}
			default:
				// Helper for close channel and set that scanner has finish
				if ofs.LastDirTime+15 < time.Now().Unix() {
					ofs.Running <- false
					close(ofs.Directories)
					return
				}
			}

		}
	}()

	go func() {
		if err := ofs.getFiles(c); err != nil {
			log.Println(err)
		}
	}()

	return ofs.FilesChan, ofs.Running, nil
}

// Token returns user token
func (ofs *OneDriveFs) Token(c client.Client) string {
	return ofs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (ofs *OneDriveFs) GetDatasourceId() string {
	return ofs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (ofs *OneDriveFs) GetThumbnail(id string, cl client.Client) (string, error) {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := fmt.Sprintf("%sitems/%s/thumbnails/0/medium", globals.OneDriveEndpoint+Drive, id)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var thumbRsp *onedrive.FileThumbnailResponse
	if err := json.NewDecoder(res.Body).Decode(&thumbRsp); err != nil {
		return "", err
	}

	return thumbRsp.URL, nil
}

// Create a one drive file
func (ofs *OneDriveFs) Create(rq file_proto.CreateRequest) chan FileMeta {
	go func() {
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
		t, err := os.Open(p)
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer t.Close()

		hc := &http.Client{}
		// https://dev.onedrive.com/items/upload_put.htm
		url := fmt.Sprintf("%sroot:/%s.%s:/content", globals.OneDriveEndpoint+Drive, rq.FileName, globals.GetDocumentTemplate(rq.MimeType, false))
		req, err := http.NewRequest("PUT", url, t) // We require a template to be able to open / edit this files online
		req.Header.Set("Authorization", ofs.token())
		req.Header.Set("Content-Type", globals.ONEDRIVE_TEXT)
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		res, err := hc.Do(req)
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer res.Body.Close()

		var f onedrive.OneDriveFile
		if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		kfo := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)

		ofs.FileMetaChan <- NewFileMeta(kfo, nil)
	}()

	return ofs.FileMetaChan
}

// Delete deletes an onedrive file
func (ofs *OneDriveFs) Delete(rq file_proto.DeleteRequest) chan FileMeta {
	go func() {
		oc := &http.Client{}
		// https://dev.onedrive.com/items/delete.htm
		url := globals.OneDriveEndpoint + Drive + "items/" + rq.OriginalId
		oreq, err := http.NewRequest("DELETE", url, nil)
		oreq.Header.Set("Authorization", ofs.token())
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		res, err := oc.Do(oreq)
		if err != nil {
			ofs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			ofs.FileMetaChan <- NewFileMeta(nil, errors.New(fmt.Sprintf("Deleting Onedrive file failed with status code %d", res.StatusCode)))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		ofs.FileMetaChan <- NewFileMeta(
			&file.KazoupOneDriveFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return ofs.FileMetaChan
}

// ShareFile
func (ofs *OneDriveFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	//POST /drive/items/{item-id}/action.invite
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	oc := &http.Client{}
	body := []byte(`{
		"requireSignIn": true,
		"sendInvitation": true,
		"roles": ["write"],
		"recipients": [
			{ "email": "` + req.DestinationId + `" }
		]
	}`)

	// https://dev.onedrive.com/items/invite.htm
	url := globals.OneDriveEndpoint + Drive + "items/" + req.OriginalId + "/action.invite"
	oreq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	oreq.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	oreq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}
	res, err := oc.Do(oreq)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return "", nil
}

// DownloadFile retrieves a file
func (ofs *OneDriveFs) DownloadFile(id string, c client.Client, opts ...string) (io.ReadCloser, error) {
	//POST /drive/items/{item-id}/action.invite
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	oc := &http.Client{}

	// https://dev.onedrive.com/items/download.htm
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/content"
	oreq, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	oreq.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	res, err := oc.Do(oreq)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// UploadFile uploads a file into google cloud storage
func (ofs *OneDriveFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, ofs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (ofs *OneDriveFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(ofs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (ofs *OneDriveFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(ofs.Endpoint.Index, "")
}

// getFiles retrieves drives, directories and files
func (ofs *OneDriveFs) getFiles(c client.Client) error {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	if err := ofs.getDrives(); err != nil {
		return err
	}
	if err := ofs.getDrivesChildren(c); err != nil {
		return err
	}

	return nil
}

// refreshToken gets a new token from custom one and saves it
func (ofs *OneDriveFs) refreshToken() error {
	tokenSource := globals.NewMicrosoftOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  ofs.Endpoint.Token.AccessToken,
		TokenType:    ofs.Endpoint.Token.TokenType,
		RefreshToken: ofs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(ofs.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return err
	}
	ofs.Endpoint.Token.AccessToken = t.AccessToken
	ofs.Endpoint.Token.TokenType = t.TokenType
	ofs.Endpoint.Token.RefreshToken = t.RefreshToken
	ofs.Endpoint.Token.Expiry = t.Expiry.Unix()

	b, err := json.Marshal(ofs.Endpoint)
	if err != nil {
		return err
	}

	c := db_proto.NewDBClient("", nil)
	_, err = c.Update(context.Background(), &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    ofs.Endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}

// token returns authorization header
func (ofs *OneDriveFs) token() string {
	return ofs.Endpoint.Token.TokenType + " " + ofs.Endpoint.Token.AccessToken
}

// getDrives retrieve user drives
func (ofs *OneDriveFs) getDrives() error {
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := globals.OneDriveEndpoint + Drives
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
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
func (ofs *OneDriveFs) getDrivesChildren(cl client.Client) error {
	var url string
	c := &http.Client{}

	for _, v := range ofs.DrivesId {
		//https://api.onedrive.com/v1.0/drives/f5a34c5d0f17415a/root/children
		url = globals.OneDriveEndpoint + Drives + v + "/root/children"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
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
				if err := ofs.pushToFilesChannel(cl, v); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// getDirChildren get children from directory
func (ofs *OneDriveFs) getDirChildren(cl client.Client, id string) error {
	// https://api.onedrive.com/v1.0/drive/items/F5A34C5D0F17415A!114/children
	c := &http.Client{}
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/children"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
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
			if err := ofs.pushToFilesChannel(cl, v); err != nil {
				return err
			}
		}
	}

	return nil
}

// pushToFilesChannel
func (ofs *OneDriveFs) pushToFilesChannel(c client.Client, f onedrive.OneDriveFile) error {
	kof := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)

	if err := ofs.generateThumbnail(c, f, kof.ID); err != nil {
		log.Println(err)
	}

	ofs.FilesChan <- kof

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (ofs *OneDriveFs) generateThumbnail(c client.Client, f onedrive.OneDriveFile, id string) error {
	n := strings.Split(f.Name, ".")

	if categories.GetDocType("."+n[len(n)-1]) == globals.CATEGORY_PICTURE {
		pr, err := ofs.DownloadFile(f.ID, c)
		if err != nil {
			return errors.New("ERROR downloading onedrive file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for onedrive file")
		}

		if err := ofs.UploadFile(b, id); err != nil {
			return errors.New("ERROR uploading thumbnail for onedrive file")
		}
	}

	return nil
}
