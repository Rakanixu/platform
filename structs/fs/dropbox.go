package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/dropbox"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"log"
	"net/http"
	"net/url"
	"os"
)

type DropboxFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

func NewDropboxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &DropboxFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

func (dfs *DropboxFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := dfs.getFiles(); err != nil {
			log.Println("ERROR geting files from dropbox ", err.Error())
		}

		dfs.Running <- false
	}()

	return dfs.FilesChan, dfs.Running, nil
}

func (dfs *DropboxFs) Token() string {
	return "Bearer " + dfs.Endpoint.Token.AccessToken
}

func (dfs *DropboxFs) GetDatasourceId() string {
	return dfs.Endpoint.Id
}

func (dfs *DropboxFs) GetThumbnail(id string) (string, error) {
	args := `{"path":"` + id + `","size":{".tag":"w640h480"}}`
	url := fmt.Sprintf("%s?authorization=%s&arg=%s", globals.DropboxThumbnailEndpoint, dfs.Token(), url.QueryEscape(args))

	return url, nil
}

func (dfs *DropboxFs) CreateFile(fileType string) (string, error) {
	// https://www.dropbox.com/developers/documentation/http/documentation#files-upload
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	p := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(fileType, true))
	t, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer t.Close()

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFileUpload, t)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", dfs.Token())
	req.Header.Set("Dropbox-API-Arg", `{
		"path": "/untitle.`+globals.GetDocumentTemplate(fileType, false)+`",
		"mode": "add",
		"autorename": true,
		"mute": false
	}`)
	req.Header.Set("Content-Type", "application/octet-stream")
	rsp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var df *dropbox.DropboxFile
	if err := json.NewDecoder(rsp.Body).Decode(&df); err != nil {
		return "", err
	}

	kfd := file.NewKazoupFileFromDropboxFile(df, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
	if err := file.IndexAsync(kfd, globals.FilesTopic, dfs.Endpoint.Index); err != nil {
		return "", err
	}

	return kfd.GetURL(), nil
}

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
