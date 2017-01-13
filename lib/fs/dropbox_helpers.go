package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"log"
	"net/http"
	"strings"
)

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
	// Retrieve all publicly available files
	// We need to retrieve like that because dropbox API does not return that information in other way
	if err := dfs.getSharedLinksForUser(""); err != nil {
		return err
	}

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

// getSharedLinksForUser retrieves all public files for a user
func (dfs *DropboxFs) getSharedLinksForUser(cursor string) error {
	b := []byte(`{}`)

	if len(cursor) > 0 {
		b = []byte(`{
			"cursor": "` + cursor + `"
		}`)
	}

	// https://www.dropbox.com/developers/documentation/http/documentation#sharing-list_shared_links
	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxSharedLinks, bytes.NewBuffer(b))
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

	var filesRsp *dropbox.PublicFilesListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&filesRsp); err != nil {
		return err
	}

	dfs.PublicFiles = append(dfs.PublicFiles, filesRsp.Links...)

	if filesRsp.HasMore {
		if err := dfs.getSharedLinksForUser(filesRsp.Cursor); err != nil {
			return err
		}
	}

	return nil
}

func (dfs *DropboxFs) pushFilesToChannel(list *dropbox.FilesListResponse) {
	var err error

	for _, v := range list.Entries {
		f := file.NewKazoupFileFromDropboxFile(v, dfs.Endpoint.Id, dfs.Endpoint.UserId, dfs.Endpoint.Index)
		if f != nil {
			f.Access = globals.ACCESS_PRIVATE

			// File is shared, lets get Users and Invitees to this file
			if f.Original.HasExplicitSharedMembers {
				f.Access = globals.ACCESS_SHARED

				f, err = dfs.getFileMembers(f)
				if err != nil {
					log.Println("ERROR getFileMembers dropbox", err)
				}
			} else {
				// File is not share, but that means to dropbox that can be private, or public (everyone with link can access the file)
				for k, v := range dfs.PublicFiles {
					// Found
					if f.Original.ID == v.ID {
						f.Access = globals.ACCESS_PUBLIC

						// Remove found for performance and break
						dfs.PublicFiles = append(dfs.PublicFiles[:k], dfs.PublicFiles[k+1:]...)
						break
					}
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
