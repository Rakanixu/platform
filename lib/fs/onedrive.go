package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/onedrive"
	"log"
	"net/http"
	"os"
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
