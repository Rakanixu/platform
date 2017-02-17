package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"log"
	"net/http"
	"os"
)

// DropboxFs dropbox file system
type DropboxFs struct {
	Endpoint            *datasource_proto.Endpoint
	PublicFiles         []dropbox.DropboxPublicFile
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

// NewDropboxFsFromEndpoint constructor
func NewDropboxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &DropboxFs{
		Endpoint:            e,
		PublicFiles:         make([]dropbox.DropboxPublicFile, 0),
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover files in dropbox datasource
func (dfs *DropboxFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := dfs.getFiles(); err != nil {
			log.Println("ERROR geting files from dropbox ", err.Error())
		}

		dfs.WalkRunning <- false
	}()

	return dfs.FilesChan, dfs.WalkRunning
}

// WalkUsers
func (dfs *DropboxFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		dfs.WalkUsersRunning <- false
	}()

	return dfs.UsersChan, dfs.WalkUsersRunning
}

// WalkChannels
func (dfs *DropboxFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		dfs.WalkChannelsRunning <- false
	}()

	return dfs.ChannelsChan, dfs.WalkChannelsRunning
}

// Enrich
func (dfs *DropboxFs) Enrich(f file.File) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupDropboxFile)
		if !ok {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		process := struct {
			Picture  bool
			Document bool
		}{
			Picture:  false,
			Document: false,
		}
		if f.(*file.KazoupDropboxFile).OptsKazoupFile == nil {
			process.Picture = true
			process.Document = true
		} else {
			process.Picture = f.(*file.KazoupDropboxFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
			process.Document = f.(*file.KazoupDropboxFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupDropboxFile).Modified)
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_PICTURE && process.Picture {
			f, err = dfs.processImage(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		if f.(*file.KazoupDropboxFile).Category == globals.CATEGORY_DOCUMENT && process.Document {
			f, err = dfs.processDocument(f.(*file.KazoupDropboxFile))
			if err != nil {
				dfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		dfs.FilesChan <- NewFileMsg(f, err)
	}()

	return dfs.FilesChan
}

// Create a file in dropbox
func (dfs *DropboxFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
		t, err := os.Open(p)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer t.Close()

		hc := &http.Client{}
		req, err := http.NewRequest("POST", globals.DropboxFileUpload, t)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		req.Header.Set("Authorization", dfs.token())
		req.Header.Set("Dropbox-API-Arg", `{"path": "/`+rq.FileName+`.`+globals.GetDocumentTemplate(rq.MimeType, false)+`","mode": "add","autorename": true,"mute": false}`)
		req.Header.Set("Content-Type", "application/octet-stream")
		rsp, err := hc.Do(req)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		var df *dropbox.DropboxFile
		if err := json.NewDecoder(rsp.Body).Decode(&df); err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		kfd := file.NewKazoupFileFromDropboxFile(*df, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		if kfd == nil {
			dfs.FilesChan <- NewFileMsg(nil, errors.New("ERROR dropbox file is nil"))
			return
		}

		dfs.FilesChan <- NewFileMsg(kfd, nil)
	}()

	return dfs.FilesChan
}

// DeleteFile deletes a dropbox file
func (dfs *DropboxFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	go func() {
		// https://www.dropbox.com/developers/documentation/http/documentation#files-delete
		b := []byte(`{
			"path": "` + rq.OriginalFilePath + `"
		}`)

		// Move file to trash in dropbox
		dc := &http.Client{}
		r, err := http.NewRequest("POST", globals.DropboxFileDelete, bytes.NewBuffer(b))
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", dfs.token())
		r.Header.Set("Content-Type", "application/json")
		rsp, err := dc.Do(r)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		defer rsp.Body.Close()

		// Check is successfully deleted
		if rsp.StatusCode != http.StatusOK {
			dfs.FilesChan <- NewFileMsg(nil, errors.New(fmt.Sprintf("Deleting Dropbox file failed with status code %d", rsp.StatusCode)))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		dfs.FilesChan <- NewFileMsg(
			&file.KazoupDropboxFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return dfs.FilesChan
}

// Update file
func (dfs *DropboxFs) Update(req file_proto.ShareRequest) chan FileMsg {
	go func() {
		// https://www.dropbox.com/developers/documentation/http/documentation#sharing-add_file_member
		// access_level cannot be editor, Dropbox API fails. Role should be selected on frontend
		b := []byte(`{
			"file": "` + req.OriginalId + `",
			"members": [
				{
					".tag": "email",
					"email": "` + req.DestinationId + `"
				}
			],
			"quiet": false,
			"access_level": "viewer",
			"add_message_as_comment": false
		}`)
		dc := &http.Client{}
		r, err := http.NewRequest("POST", globals.DropboxFileShare, bytes.NewBuffer(b))
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", dfs.token())
		r.Header.Set("Content-Type", "application/json")
		rsp, err := dc.Do(r)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusOK {
			dfs.FilesChan <- NewFileMsg(nil, errors.New(fmt.Sprintf("Sharing Dropbox file failed with status code %d", rsp.StatusCode)))
			return
		}

		// Get the modified file to reindex
		f, err := dfs.getFile(req.OriginalId)
		if err != nil {
			dfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		dfs.FilesChan <- NewFileMsg(f, nil)
	}()

	return dfs.FilesChan
}
