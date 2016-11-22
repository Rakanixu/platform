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
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// DropboxFs dropbox file system
type DropboxFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

// NewDropboxFsFromEndpoint constructor
func NewDropboxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &DropboxFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

// List returns 2 channels, for files and state. Discover files in dropbox datasource
func (dfs *DropboxFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := dfs.getFiles(); err != nil {
			log.Println("ERROR geting files from dropbox ", err.Error())
		}

		dfs.Running <- false
	}()

	return dfs.FilesChan, dfs.Running, nil
}

// Token returns dropbox user token
func (dfs *DropboxFs) Token() string {
	return "Bearer " + dfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (dfs *DropboxFs) GetDatasourceId() string {
	return dfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (dfs *DropboxFs) GetThumbnail(id string) (string, error) {
	args := `{"path":"` + id + `","size":{".tag":"w640h480"}}`
	url := fmt.Sprintf("%s?authorization=%s&arg=%s", globals.DropboxThumbnailEndpoint, dfs.Token(), url.QueryEscape(args))

	return url, nil
}

// CreateFile creates a file in dropbox and index it on Elastic Search
func (dfs *DropboxFs) CreateFile(rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
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

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFileUpload, t)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.Token())
	req.Header.Set("Dropbox-API-Arg", `{
		"path": "/`+rq.FileName+`.`+globals.GetDocumentTemplate(rq.MimeType, false)+`",
		"mode": "add",
		"autorename": true,
		"mute": false
	}`)
	req.Header.Set("Content-Type", "application/octet-stream")
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var df *dropbox.DropboxFile
	if err := json.NewDecoder(rsp.Body).Decode(&df); err != nil {
		return nil, err
	}

	kfd := file.NewKazoupFileFromDropboxFile(df, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
	if err := file.IndexAsync(kfd, globals.FilesTopic, dfs.Endpoint.Index, true); err != nil {
		return nil, err
	}

	b, err := json.Marshal(kfd)
	if err != nil {
		return nil, err
	}

	return &file_proto.CreateResponse{
		DocUrl: kfd.GetURL(),
		Data:   string(b),
	}, nil
}

// DeleteFile deletes a dropbox file
func (dfs *DropboxFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#files-delete
	b := []byte(`{
		"path": "` + rq.OriginalFilePath + `"
	}`)

	// Move file to trash in dropbox
	dc := &http.Client{}
	r, err := http.NewRequest("POST", globals.DropboxFileDelete, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", dfs.Token())
	r.Header.Set("Content-Type", "application/json")
	rsp, err := dc.Do(r)
	if err != nil {
		return nil, err
	}
	// From Dropbox docs:
	// The returned metadata will be the corresponding FileMetadata or FolderMetadata for the item
	// at time of deletion, and not a DeletedMetadata object.
	// Obviously, we need the DeletedMetadata, so we have to get it by doing fucking extra request
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Deleting Dropbox file failed with status code %d", rsp.StatusCode))
	}

	// Reindex trashed file
	// TODO: this fails and do not see why it should
	/*	kfd, err := dfs.getFile(rq.OriginalId)
		if err != nil {
			return nil, err
		}

		log.Println(kfd)
		//log.Println(df.Tag)

		if err := file.IndexAsync(kfd, globals.FilesTopic, dfs.Endpoint.Index); err != nil {
			return nil, err
		}*/

	// HACK: read the file and update manually, well, it works, but timestamps are old ones
	// FIXME: Try to debug commented code because it should work.
	rreq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: rq.Index,
			Type:  globals.FileType,
			Id:    rq.FileId,
		},
	)
	rres := &db_proto.ReadResponse{}
	if err := c.Call(ctx, rreq, rres); err != nil {
		return nil, err
	}

	var df *file.KazoupDropboxFile
	if err := json.Unmarshal([]byte(rres.Result), &df); err != nil {
		return nil, err
	}

	// HACK HACK HACK!
	df.Original.DropboxTag = "deleted"

	bb, err := json.Marshal(df)
	if err != nil {
		return nil, err
	}

	ureq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Update",
		&db_proto.UpdateRequest{
			Index: rq.Index,
			Type:  globals.FileType,
			Id:    rq.FileId,
			Data:  string(bb),
		},
	)
	ures := &db_proto.UpdateResponse{}
	if err := c.Call(ctx, ureq, ures); err != nil {
		return nil, err
	}

	return &file_proto.DeleteResponse{}, nil
}

// ShareFile
func (dfs *DropboxFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
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
		return "", err
	}
	r.Header.Set("Authorization", dfs.Token())
	r.Header.Set("Content-Type", "application/json")
	rsp, err := dc.Do(r)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Sharing Dropbox file failed with status code %d", rsp.StatusCode))
	}

	// Get the modified file to reindex
	f, err := dfs.getFile(req.OriginalId)
	if err != nil {
		return "", err
	}

	if err := file.IndexAsync(f, globals.FilesTopic, dfs.Endpoint.Index, true); err != nil {
		return "", err
	}

	return "", nil
}

// DownloadFile retrieves a file
func (dfs *DropboxFs) DownloadFile(id string, opts ...string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, globals.DropboxFileDownload, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.Token())
	req.Header.Set("Dropbox-API-Arg", `{
		"path": "`+id+`"
	}`)
	rsp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UploadFile uploads a file into google cloud storage
func (dfs *DropboxFs) UploadFile(file []byte, fId string) error {
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
	r.Header.Set("Authorization", dfs.Token())
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

	kfd := file.NewKazoupFileFromDropboxFile(f, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)

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
	req.Header.Set("Authorization", dfs.Token())
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

	for _, v := range filesRsp.Entries {
		name := strings.Split(v.Name, ".")
		if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
			b, err := dfs.DownloadFile(v.ID)
			if err != nil {
				log.Println("ERROR downloading dropbox file: %s", err)
			}

			b, err = image.Thumbnail(b, globals.THUMBNAIL_WIDTH)
			if err != nil {
				log.Println("ERROR generating thumbnail for dropbox file: %s", err)
			}

			if err := dfs.UploadFile(b, v.ID); err != nil {
				log.Println("ERROR uploading thumbnail for dropbox file: %s", err)
			}
		}

		f := file.NewKazoupFileFromDropboxFile(&v, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		// File is shared, lets get Users and Invitees to this file
		if f.Original.HasExplicitSharedMembers {
			f, err = dfs.getFileMembers(f)
			if err != nil {
				return err
			}
		}

		dfs.FilesChan <- f
	}

	if filesRsp.HasMore {
		dfs.getNextPage(filesRsp.Cursor)
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
	req.Header.Set("Authorization", dfs.Token())
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

	for _, v := range filesRsp.Entries {
		f := file.NewKazoupFileFromDropboxFile(&v, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		dfs.FilesChan <- f
	}

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
	req.Header.Set("Authorization", dfs.Token())
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
	req.Header.Set("Authorization", dfs.Token())
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
