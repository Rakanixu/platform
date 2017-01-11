package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/onedrive"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
	DrivesId            []string
	LastDirTime         int64
	Directories         chan string
}

// NewOneDriveFsFromEndpoint constructor
func NewOneDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &OneDriveFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
		DrivesId:            []string{},
		// This is important to have a size bigger than one, the bigger, less likely to block
		// If not, program execution will block, due to recursivity,
		// We are pushing more elements before finish execution.
		// I expect to never push 10000 folders before other folders have been completly scanned
		Directories: make(chan string, 10000),
	}
}

// Walk returns 2 channels, for files and state. Discover files in one drive datasources
func (ofs *OneDriveFs) Walk() (chan FileMsg, chan bool) {
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
					ofs.WalkRunning <- false
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

	return ofs.FilesChan, ofs.WalkRunning
}

// WalkUsers
func (ofs *OneDriveFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		ofs.WalkUsersRunning <- false
	}()

	return ofs.UsersChan, ofs.WalkUsersRunning
}

// WalkChannels
func (ofs *OneDriveFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		ofs.WalkChannelsRunning <- false
	}()

	return ofs.ChannelsChan, ofs.WalkChannelsRunning
}

// Create a one drive file
func (ofs *OneDriveFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
		t, err := os.Open(p)
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
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
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		res, err := hc.Do(req)
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer res.Body.Close()

		var f onedrive.OneDriveFile
		if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		kfo := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)

		ofs.FilesChan <- NewFileMsg(kfo, nil)
	}()

	return ofs.FilesChan
}

// Delete deletes an onedrive file
func (ofs *OneDriveFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	go func() {
		oc := &http.Client{}
		// https://dev.onedrive.com/items/delete.htm
		url := globals.OneDriveEndpoint + Drive + "items/" + rq.OriginalId
		oreq, err := http.NewRequest("DELETE", url, nil)
		oreq.Header.Set("Authorization", ofs.token())
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		res, err := oc.Do(oreq)
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusNoContent {
			ofs.FilesChan <- NewFileMsg(nil, errors.New(fmt.Sprintf("Deleting Onedrive file failed with status code %d", res.StatusCode)))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		ofs.FilesChan <- NewFileMsg(
			&file.KazoupOneDriveFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return ofs.FilesChan
}

// Update file
func (ofs *OneDriveFs) Update(req file_proto.ShareRequest) chan FileMsg {
	/*	go func() {
		//POST /drive/items/{item-id}/action.invite
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
		oreq.Header.Set("Authorization", ofs.token())
		oreq.Header.Set("Content-Type", "application/json")
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		res, err := oc.Do(oreq)
		if err != nil {
			ofs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer res.Body.Close()

		ofs.FilesChan <- NewFileMsg(
			&file.KazoupOneDriveFile{
				file.KazoupFile{
					ID: req.FileId,
				},
				nil,
			},
			nil,
		)

		// TODO: request for file to Onedrive and retur whole file
	}()*/

	return ofs.FilesChan
}

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
func (ofs *OneDriveFs) getDrivesChildren() error {
	var url string
	c := &http.Client{}

	for _, v := range ofs.DrivesId {
		//https://api.onedrive.com/v1.0/drives/f5a34c5d0f17415a/root/children
		url = globals.OneDriveEndpoint + Drives + v + "/root/children"

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
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/children"
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

// pushToFilesChannel
func (ofs *OneDriveFs) pushToFilesChannel(f onedrive.OneDriveFile) error {
	kof := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)

	if err := ofs.generateThumbnail(f, kof.ID); err != nil {
		log.Println(err)
	}

	ofs.FilesChan <- NewFileMsg(kof, nil)

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (ofs *OneDriveFs) generateThumbnail(f onedrive.OneDriveFile, id string) error {
	n := strings.Split(f.Name, ".")

	if categories.GetDocType("."+n[len(n)-1]) == globals.CATEGORY_PICTURE {
		// Download file from OneDrive, so connector is globals.OneDrive
		ocs, err := cs.NewCloudStorageFromEndpoint(ofs.Endpoint, globals.OneDrive)
		if err != nil {
			return err
		}

		pr, err := ocs.Download(f.ID)
		if err != nil {
			return errors.New("ERROR downloading onedrive file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for onedrive file")
		}
		// Upload file to GoogleCloudStorage, so connector is globals.GoogleCloudStorage
		ncs, err := cs.NewCloudStorageFromEndpoint(ofs.Endpoint, globals.GoogleCloudStorage)
		if err != nil {
			return err
		}

		if err := ncs.Upload(b, id); err != nil {
			return err
		}
	}

	return nil
}
