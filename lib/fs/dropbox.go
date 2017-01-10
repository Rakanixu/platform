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
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/micro/go-micro/client"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// DropboxFs dropbox file system
type DropboxFs struct {
	Endpoint     *datasource_proto.Endpoint
	Running      chan bool
	FilesChan    chan file.File
	FileMetaChan chan FileMeta
}

// NewDropboxFsFromEndpoint constructor
func NewDropboxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &DropboxFs{
		Endpoint:     e,
		Running:      make(chan bool, 1),
		FilesChan:    make(chan file.File),
		FileMetaChan: make(chan FileMeta),
	}
}

// List returns 2 channels, for files and state. Discover files in dropbox datasource
func (dfs *DropboxFs) List(c client.Client) (chan file.File, chan bool, error) {
	go func() {
		if err := dfs.getFiles(c); err != nil {
			log.Println("ERROR geting files from dropbox ", err.Error())
		}

		dfs.Running <- false
	}()

	return dfs.FilesChan, dfs.Running, nil
}

// Token returns dropbox user token
func (dfs *DropboxFs) Token(c client.Client) string {
	return "Bearer " + dfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (dfs *DropboxFs) GetDatasourceId() string {
	return dfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (dfs *DropboxFs) GetThumbnail(id string, c client.Client) (string, error) {
	args := `{"path":"` + id + `","size":{".tag":"w640h480"}}`
	url := fmt.Sprintf("%s?authorization=%s&arg=%s", globals.DropboxThumbnailEndpoint, dfs.Token(c), url.QueryEscape(args))

	return url, nil
}

// Create a file in dropbox
func (dfs *DropboxFs) Create(rq file_proto.CreateRequest) chan FileMeta {
	go func() {
		// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
		t, err := os.Open(p)
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer t.Close()

		hc := &http.Client{}
		req, err := http.NewRequest("POST", globals.DropboxFileUpload, t)
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
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
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer rsp.Body.Close()

		var df *dropbox.DropboxFile
		if err := json.NewDecoder(rsp.Body).Decode(&df); err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		kfd := file.NewKazoupFileFromDropboxFile(*df, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		if kfd == nil {
			dfs.FileMetaChan <- NewFileMeta(nil, errors.New("ERROR dropbox file is nil"))
			return
		}

		dfs.FileMetaChan <- NewFileMeta(kfd, nil)
	}()

	return dfs.FileMetaChan
}

// DeleteFile deletes a dropbox file
func (dfs *DropboxFs) Delete(rq file_proto.DeleteRequest) chan FileMeta {
	go func() {
		// https://www.dropbox.com/developers/documentation/http/documentation#files-delete
		b := []byte(`{
			"path": "` + rq.OriginalFilePath + `"
		}`)

		// Move file to trash in dropbox
		dc := &http.Client{}
		r, err := http.NewRequest("POST", globals.DropboxFileDelete, bytes.NewBuffer(b))
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		r.Header.Set("Authorization", dfs.token())
		r.Header.Set("Content-Type", "application/json")
		rsp, err := dc.Do(r)
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		defer rsp.Body.Close()

		// Check is successfully deleted
		if rsp.StatusCode != http.StatusOK {
			dfs.FileMetaChan <- NewFileMeta(nil, errors.New(fmt.Sprintf("Deleting Dropbox file failed with status code %d", rsp.StatusCode)))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		dfs.FileMetaChan <- NewFileMeta(
			&file.KazoupDropboxFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return dfs.FileMetaChan
}

// Update file
func (dfs *DropboxFs) Update(req file_proto.ShareRequest) chan FileMeta {
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
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		r.Header.Set("Authorization", dfs.token())
		r.Header.Set("Content-Type", "application/json")
		rsp, err := dc.Do(r)
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusOK {
			dfs.FileMetaChan <- NewFileMeta(nil, errors.New(fmt.Sprintf("Sharing Dropbox file failed with status code %d", rsp.StatusCode)))
			return
		}

		// Get the modified file to reindex
		f, err := dfs.getFile(req.OriginalId)
		if err != nil {
			dfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		dfs.FileMetaChan <- NewFileMeta(f, nil)
	}()

	return dfs.FileMetaChan
}

// DownloadFile retrieves a file
func (dfs *DropboxFs) DownloadFile(id string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, globals.DropboxFileDownload, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Dropbox-API-Arg", `{
			"path": "`+id+`"
		}`)
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

// UploadFile uploads a file into google cloud storage
func (dfs *DropboxFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, dfs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (dfs *DropboxFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(dfs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (dfs *DropboxFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(dfs.Endpoint.Index, "")
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
func (dfs *DropboxFs) getFiles(cl client.Client) error {
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
	req.Header.Set("Authorization", dfs.Token(cl))
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
		dfs.getNextPage(filesRsp.Cursor, cl)
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (dfs *DropboxFs) generateThumbnail(f dropbox.DropboxFile, id string) error {
	name := strings.Split(f.Name, ".")

	if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
		pr, err := dfs.DownloadFile(f.ID)
		if err != nil {
			return errors.New("ERROR downloading dropbox file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for dropbox file")
		}

		if err := dfs.UploadFile(b, id); err != nil {
			return errors.New("ERROR uploading thumbnail for dropbox file")
		}
	}

	return nil
}

// getNextPage allows pagination while discovering files
func (dfs *DropboxFs) getNextPage(cursor string, cl client.Client) error {
	// https://www.dropbox.com/developers/documentation/http/documentation#files-list_folder-continue
	b := []byte(`{
		"cursor":"` + cursor + `"
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFilesEndpoint+"/continue", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", dfs.Token(cl))
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
		dfs.getNextPage(filesRsp.Cursor, cl)
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

			dfs.FilesChan <- f
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
