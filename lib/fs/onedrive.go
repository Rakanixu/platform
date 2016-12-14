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
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
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
	Endpoint    *datasource_proto.Endpoint
	Running     chan bool
	FilesChan   chan file.File
	DrivesId    []string
	Directories chan string
	LastDirTime int64
}

// NewOneDriveFsFromEndpoint constructor
func NewOneDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &OneDriveFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
		DrivesId:  []string{},
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

				err := ofs.getDirChildren(v)
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
		if err := ofs.getFiles(); err != nil {
			log.Println(err)
		}
	}()

	return ofs.FilesChan, ofs.Running, nil
}

// Token returns user token
func (ofs *OneDriveFs) Token() string {
	return ofs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (ofs *OneDriveFs) GetDatasourceId() string {
	return ofs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (ofs *OneDriveFs) GetThumbnail(id string) (string, error) {
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

// CreateFile creates a one drive document and index it on Elastic Search
func (ofs *OneDriveFs) CreateFile(ctx context.Context, c client.Client, rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return nil, err
	}

	p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
	t, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer t.Close()

	hc := &http.Client{}
	// https://dev.onedrive.com/items/upload_put.htm
	url := fmt.Sprintf("%sroot:/%s.%s:/content", globals.OneDriveEndpoint+Drive, rq.FileName, globals.GetDocumentTemplate(rq.MimeType, false))
	req, err := http.NewRequest("PUT", url, t) // We require a template to be able to open / edit this files online
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	req.Header.Set("Content-Type", globals.ONEDRIVE_TEXT)
	if err != nil {
		return nil, err
	}
	res, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var f *onedrive.OneDriveFile
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		return nil, err
	}

	kfo := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)
	if err := file.IndexAsync(c, kfo, globals.FilesTopic, ofs.Endpoint.Index, true); err != nil {
		return nil, err
	}

	b, err := json.Marshal(kfo)
	if err != nil {
		return nil, err
	}

	return &file_proto.CreateResponse{
		DocUrl: kfo.GetURL(),
		Data:   string(b),
	}, nil
}

// DeleteFile deletes an onedrive file
func (ofs *OneDriveFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	oc := &http.Client{}
	// https://dev.onedrive.com/items/delete.htm
	url := globals.OneDriveEndpoint + Drive + "items/" + rq.OriginalId
	oreq, err := http.NewRequest("DELETE", url, nil)
	oreq.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	if err != nil {
		return nil, err
	}
	res, err := oc.Do(oreq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return nil, errors.New(fmt.Sprintf("Deleting Onedrive file failed with status code %d", res.StatusCode))
	}

	// Remove file from index
	req := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Delete",
		&db_proto.DeleteRequest{
			Index: rq.Index,
			Type:  globals.FileType,
			Id:    rq.FileId,
		},
	)
	rsp := &db_proto.DeleteResponse{}
	if err := c.Call(ctx, req, rsp); err != nil {
		return nil, err
	}

	// Publish notification topic, let client know when to refresh itself
	if err := c.Publish(globals.NewSystemContext(), c.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
		Method: globals.NOTIFY_REFRESH_SEARCH,
		UserId: ofs.Endpoint.UserId,
	})); err != nil {
		log.Print("Publishing (notify file) error %s", err)
	}

	return &file_proto.DeleteResponse{}, nil
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
func (ofs *OneDriveFs) DownloadFile(id string, opts ...string) (io.ReadCloser, error) {
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
func (ofs *OneDriveFs) getFiles() error {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	if err := ofs.getDrives(); err != nil {
		return err
	}
	if err := ofs.getDrivesChildren(); err != nil {
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
func (ofs *OneDriveFs) getDrivesChildren() error {
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
			if err := ofs.pushToFilesChannel(v); err != nil {
				return err
			}
		}
	}

	return nil
}

// pushToFilesChannel
func (ofs *OneDriveFs) pushToFilesChannel(f onedrive.OneDriveFile) error {
	if err := ofs.generateThumbnail(f); err != nil {
		log.Println(err)
	}

	kof := file.NewKazoupFileFromOneDriveFile(&f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)
	ofs.FilesChan <- kof

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (ofs *OneDriveFs) generateThumbnail(f onedrive.OneDriveFile) error {
	n := strings.Split(f.Name, ".")

	if categories.GetDocType("."+n[len(n)-1]) == globals.CATEGORY_PICTURE {
		pr, err := ofs.DownloadFile(f.ID)
		if err != nil {
			return errors.New("ERROR downloading onedrive file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for onedrive file")
		}

		if err := ofs.UploadFile(b, f.ID); err != nil {
			return errors.New("ERROR uploading thumbnail for onedrive file")
		}
	}

	return nil
}
