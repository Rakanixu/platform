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
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"io"
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

// CreateFile creates a file in dropbox and index it on Elastic Search
func (dfs *DropboxFs) CreateFile(ctx context.Context, c client.Client, rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
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

	hc := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFileUpload, t)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.Token(c))
	req.Header.Set("Dropbox-API-Arg", `{
		"path": "/`+rq.FileName+`.`+globals.GetDocumentTemplate(rq.MimeType, false)+`",
		"mode": "add",
		"autorename": true,
		"mute": false
	}`)
	req.Header.Set("Content-Type", "application/octet-stream")
	rsp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var df *dropbox.DropboxFile
	if err := json.NewDecoder(rsp.Body).Decode(&df); err != nil {
		return nil, err
	}

	kfd := file.NewKazoupFileFromDropboxFile(*df, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
	if kfd == nil {
		return nil, errors.New("ERROR dropbox file is nil")
	}

	if err := file.IndexAsync(c, kfd, globals.FilesTopic, dfs.Endpoint.Index, true); err != nil {
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
	r.Header.Set("Authorization", dfs.Token(c))
	r.Header.Set("Content-Type", "application/json")
	rsp, err := dc.Do(r)
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	// Check is successfully deleted
	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Deleting Dropbox file failed with status code %d", rsp.StatusCode))
	}

	// Delete from index
	req := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Delete",
		&db_proto.DeleteRequest{
			Index: dfs.Endpoint.Index,
			Type:  globals.FileType,
			Id:    rq.FileId,
		},
	)
	res := &db_proto.DeleteResponse{}

	if err := c.Call(ctx, req, res); err != nil {
		return nil, err
	}

	// Publish notification topic, let client know when to refresh itself
	if err := c.Publish(globals.NewSystemContext(), c.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
		Method: globals.NOTIFY_REFRESH_SEARCH,
		UserId: dfs.Endpoint.UserId,
	})); err != nil {
		log.Print("Publishing (notify file) error %s", err)
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
	r.Header.Set("Authorization", dfs.Token(c))
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
	f, err := dfs.getFile(req.OriginalId, c)
	if err != nil {
		return "", err
	}

	if err := file.IndexAsync(c, f, globals.FilesTopic, dfs.Endpoint.Index, true); err != nil {
		return "", err
	}

	return "", nil
}

// DownloadFile retrieves a file
func (dfs *DropboxFs) DownloadFile(id string, cl client.Client, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, globals.DropboxFileDownload, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.Token(cl))
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
func (dfs *DropboxFs) getFile(id string, c client.Client) (*file.KazoupDropboxFile, error) {
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
	r.Header.Set("Authorization", dfs.Token(c))
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

	return dfs.getFileMembers(kfd, c)
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

	dfs.pushFilesToChannel(filesRsp, cl)

	if filesRsp.HasMore {
		dfs.getNextPage(filesRsp.Cursor, cl)
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (dfs *DropboxFs) generateThumbnail(f dropbox.DropboxFile, id string, c client.Client) error {
	name := strings.Split(f.Name, ".")

	if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
		pr, err := dfs.DownloadFile(f.ID, c)
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

	dfs.pushFilesToChannel(filesRsp, cl)

	if filesRsp.HasMore {
		dfs.getNextPage(filesRsp.Cursor, cl)
	}

	return nil
}

// getFileMembers retrieves users with acces to a given file
func (dfs *DropboxFs) getFileMembers(f *file.KazoupDropboxFile, cl client.Client) (*file.KazoupDropboxFile, error) {
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
	req.Header.Set("Authorization", dfs.Token(cl))
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
			a, err := dfs.getAccount(v.User.AccountID, cl)
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

func (dfs *DropboxFs) pushFilesToChannel(list *dropbox.FilesListResponse, c client.Client) {
	var err error

	for _, v := range list.Entries {
		f := file.NewKazoupFileFromDropboxFile(v, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		if f != nil {
			// File is shared, lets get Users and Invitees to this file
			if f.Original.HasExplicitSharedMembers {
				f, err = dfs.getFileMembers(f, c)
				if err != nil {
					log.Println("ERROR getFileMembers dropbox", err)
				}
			}

			if err := dfs.generateThumbnail(v, f.ID, c); err != nil {
				log.Println(err)
			}

			dfs.FilesChan <- f
		}
	}
}

// getAccount retrieves dropbox user accounts
func (dfs *DropboxFs) getAccount(aId string, cl client.Client) (*dropbox.DropboxUser, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#users-get_account
	b := []byte(`{
		"account_id":"` + aId + `"
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxAccountEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", dfs.Token(cl))
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
