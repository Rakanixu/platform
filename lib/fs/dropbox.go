package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"log"
	"net/http"
	"os"
	"strings"
)

// DropboxFs dropbox file system
type DropboxFs struct {
	Endpoint            *datasource_proto.Endpoint
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
		req.Header.Set("Dropbox-API-Arg", `{
			"path": "/`+rq.FileName+`.`+globals.GetDocumentTemplate(rq.MimeType, false)+`",
			"mode": "add",
			"autorename": true,
			"mute": false
		}`)
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

// getFile retrieves a single file from dorpbox
func (dfs *DropboxFs) getFile(id string) (*file.KazoupDropboxFile, error) {
	b := []byte(`{
		"path": "` + id + `",
		"include_media_info": true,
		"include_deleted": true,
		"include_has_explicit_shared_members": true
	}`)

	dc := &http.Client{}
	r, err := http.NewRequest("POST", globals.DropboxFileEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", dfs.token())
	r.Header.Set("Content-Type", "application/json")
	rsp, err := dc.Do(r)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var f *dropbox.DropboxFile
	if err := json.NewDecoder(rsp.Body).Decode(&f); err != nil {
		return nil, err
	}

	kfd := file.NewKazoupFileFromDropboxFile(*f, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
	if kfd == nil {
		return nil, errors.New("ERROR dropbox file is nil")
	}

	return dfs.getFileMembers(kfd)
}

// getFiles discovers files in dropbox account
func (dfs *DropboxFs) getFiles() error {
	// We want all avilable info
	// https://dropbox.github.io/dropbox-api-v2-explorer/#files_list_folder
	b := []byte(`{
		"path":"",
		"recursive":true,
		"include_media_info":true,
		"include_deleted":true,
		"include_has_explicit_shared_members":true
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFilesEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var filesRsp *dropbox.FilesListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&filesRsp); err != nil {
		return err
	}

	dfs.pushFilesToChannel(filesRsp)

	if filesRsp.HasMore {
		dfs.getNextPage(filesRsp.Cursor)
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (dfs *DropboxFs) generateThumbnail(f dropbox.DropboxFile, id string) error {
	name := strings.Split(f.Name, ".")

	if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
		// Downloads from dropbox, see connector
		dcs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.Dropbox)
		if err != nil {
			return err
		}

		pr, err := dcs.Download(f.ID)
		if err != nil {
			return errors.New("ERROR downloading dropbox file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for dropbox file")
		}

		// Uploads to Google cloud storage, see connector
		ncs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.GoogleCloudStorage)
		if err != nil {
			return err
		}

		if err := ncs.Upload(b, id); err != nil {
			return err
		}
	}

	return nil
}

// getNextPage allows pagination while discovering files
func (dfs *DropboxFs) getNextPage(cursor string) error {
	// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder-continue
	b := []byte(`{
		"cursor":"` + cursor + `"
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFilesEndpoint+"/continue", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var filesRsp *dropbox.FilesListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&filesRsp); err != nil {
		return err
	}

	dfs.pushFilesToChannel(filesRsp)

	if filesRsp.HasMore {
		dfs.getNextPage(filesRsp.Cursor)
	}

	return nil
}

// getFileMembers retrieves users with acces to a given file
func (dfs *DropboxFs) getFileMembers(f *file.KazoupDropboxFile) (*file.KazoupDropboxFile, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#sharing-list_file_members
	b := []byte(`{
		"file":"` + f.Original.ID + `",
		"include_inherited": true,
		"limit": 250
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFileMembers, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var membersRsp *dropbox.FileMembersListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&membersRsp); err != nil {
		return nil, err
	}

	if len(membersRsp.Users) > 0 {
		f.Original.DropboxUsers = make([]dropbox.DropboxUser, 0)

		for _, v := range membersRsp.Users {
			a, err := dfs.getAccount(v.User.AccountID)
			if err != nil {
				return nil, err
			}

			f.Original.DropboxUsers = append(f.Original.DropboxUsers, *a)
		}
	}

	if len(membersRsp.Invitees) > 0 {
		f.Original.DropboxInvitees = membersRsp.Invitees
	}

	// TODO: membersRsp.Groups, I just ignore, we can attach them to the DropboxFile, so can be used in front

	return f, nil
}

func (dfs *DropboxFs) pushFilesToChannel(list *dropbox.FilesListResponse) {
	var err error

	for _, v := range list.Entries {
		f := file.NewKazoupFileFromDropboxFile(v, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		if f != nil {
			// File is shared, lets get Users and Invitees to this file
			if f.Original.HasExplicitSharedMembers {
				f, err = dfs.getFileMembers(f)
				if err != nil {
					log.Println("ERROR getFileMembers dropbox", err)
				}
			}

			if err := dfs.generateThumbnail(v, f.ID); err != nil {
				log.Println(err)
			}

			dfs.FilesChan <- NewFileMsg(f, nil)
		}
	}
}

// getAccount retrieves dropbox user accounts
func (dfs *DropboxFs) getAccount(aId string) (*dropbox.DropboxUser, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#users-get_account
	b := []byte(`{
		"account_id":"` + aId + `"
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxAccountEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var account *dropbox.DropboxUser
	if err := json.NewDecoder(rsp.Body).Decode(&account); err != nil {
		return nil, err
	}

	return account, nil
}

// token returns auth header token
func (dfs *DropboxFs) token() string {
	return "Bearer " + dfs.Endpoint.Token.AccessToken
}
