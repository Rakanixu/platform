package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// BoxFs Box File System
type BoxFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
	Directories         chan string
	LastDirTime         int64
	DefaultOffset       int
	DefaultLimit        int
}

// NewBoxFsFromEndpoint constructor
func NewBoxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &BoxFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
		// This is important to have a size bigger than one, the bigger, less likely to block
		// If not, program execution will block, due to recursivity,
		// We are pushing more elements before finish execution.
		// I expect to never push 10000 folders before other folders have been completly scanned
		Directories:   make(chan string, 10000),
		DefaultOffset: 0,
		DefaultLimit:  100,
	}
}

// List returns 2 channels, one for files , other for the state. Goes over a datasource and discover files
func (bfs *BoxFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		bfs.LastDirTime = time.Now().Unix()
		for {
			select {
			case v := <-bfs.Directories:
				bfs.LastDirTime = time.Now().Unix()

				err := bfs.getDirChildren(v, bfs.DefaultOffset, bfs.DefaultLimit)
				if err != nil {
					log.Println(err)
				}
			default:
				// Helper for close channel and set that scanner has finish
				if bfs.LastDirTime+10 < time.Now().Unix() {
					close(bfs.Directories)
					bfs.WalkRunning <- false
					return
				}
			}

		}
	}()

	go func() {
		if err := bfs.getDirChildren("0", bfs.DefaultOffset, bfs.DefaultLimit); err != nil {
			log.Println(err)
		}
	}()

	return bfs.FilesChan, bfs.WalkRunning
}

// WalkUsers
func (bfs *BoxFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		bfs.WalkUsersRunning <- false
	}()

	return bfs.UsersChan, bfs.WalkUsersRunning
}

// WalkChannels
func (bfs *BoxFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		bfs.WalkChannelsRunning <- false
	}()

	return bfs.ChannelsChan, bfs.WalkChannelsRunning
}

// Create file in box
func (bfs *BoxFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		// Box supports multi part form upload
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		// File template path
		t := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", utils.GetDocumentTemplate(rq.MimeType, true))
		buf := &bytes.Buffer{}
		mw := multipart.NewWriter(buf)
		defer mw.Close()

		f, err := os.Open(t)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer f.Close()

		// This is how you upload a file as multipart form
		ff, err := mw.CreateFormFile("file", t)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		if _, err = io.Copy(ff, f); err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		// Add extra fields required by API
		mw.WriteField(
			"attributes",
			`{"name":"`+rq.FileName+`.`+utils.GetDocumentTemplate(rq.MimeType, false)+`", "parent":{"id":"0"}}`,
		)
		if err := mw.Close(); err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		hc := &http.Client{}
		req, err := http.NewRequest("POST", globals.BoxUploadEndpoint, buf)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		req.Header.Set("Authorization", bfs.token())
		req.Header.Set("Content-Type", mw.FormDataContentType())

		rsp, err := hc.Do(req)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		var bf *box.BoxUpload
		if err := json.NewDecoder(rsp.Body).Decode(&bf); err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		if rsp.StatusCode == http.StatusConflict {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Conflict creating file in Box, file with same name already exists"))
			return
		}

		if rsp.StatusCode != http.StatusCreated && bf.TotalCount != 1 {
			bfs.FilesChan <- NewFileMsg(nil, errors.New("Failed creating file in Box"))
			return
		}

		// Construct Kazoup file from box created file and index it
		kfb := file.NewKazoupFileFromBoxFile(bf.Entries[0], bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

		bfs.FilesChan <- NewFileMsg(kfb, nil)
	}()

	return bfs.FilesChan
}

// DeleteFile deletes a box file
func (bfs *BoxFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	go func() {
		// https://docs.box.com/reference#delete-a-file
		// Depending on the enterprise settings for this user, the item will either be actually deleted from Box or moved to the trash.
		bc := &http.Client{}
		url := fmt.Sprintf("%s%s", globals.BoxFileMetadataEndpoint, rq.OriginalId)
		r, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", bfs.token())
		rsp, err := bc.Do(r)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusNoContent {
			bfs.FilesChan <- NewFileMsg(nil, fmt.Errorf("Deleting Box file failed with status code %d", rsp.StatusCode))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		bfs.FilesChan <- NewFileMsg(
			&file.KazoupBoxFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return bfs.FilesChan
}

// Update file
func (bfs *BoxFs) Update(req file_proto.ShareRequest) chan FileMsg {
	go func() {
		b := []byte(`{
			"shared_link": {
				"access": "open"
			}
		}`)

		bc := &http.Client{}
		url := fmt.Sprintf("%s%s", globals.BoxFileMetadataEndpoint, req.OriginalId)
		r, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", bfs.token())
		rsp, err := bc.Do(r)
		if err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusOK {
			bfs.FilesChan <- NewFileMsg(nil, fmt.Errorf("Sharing Box file failed with status code %d", rsp.StatusCode))
			return
		}

		var f *box.BoxFileMeta
		if err := json.NewDecoder(rsp.Body).Decode(&f); err != nil {
			bfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		kbf := file.NewKazoupFileFromBoxFile(*f, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

		bfs.FilesChan <- NewFileMsg(kbf, nil)
	}()

	return bfs.FilesChan
}
