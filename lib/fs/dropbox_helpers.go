package fs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/objectstorage"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"io"
	"io/ioutil"
	"net/http"
	"time"
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

// processImage, thumbnail generation, cloud vision processing
func (dfs *DropboxFs) processImage(f *file.KazoupDropboxFile) (file.File, error) {
	// Downloads from dropbox, see connector
	dcs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.Dropbox)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err = backoff.Retry(func() error {
		rc, err = dcs.Download(f.OriginalID)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	// Resize to optimal size for cloud vision API
	cvrd, err := image.Thumbnail(rc, globals.CLOUD_VISION_IMG_WIDTH)
	if err != nil {
		return nil, err
	}

	// Library implements exponential backoff
	if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
		return nil, err
	}

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			TagsTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.TagsTimestamp = &n
	}

	return f, nil
}

// enrichFile sends the original file to tika and enrich KazoupDropboxFile with Tika interface
func (dfs *DropboxFs) processDocument(f *file.KazoupDropboxFile) (file.File, error) {
	// Download file from Dropbox, so connector is globals.Dropbox
	dcs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.Dropbox)
	if err != nil {
		return nil, err
	}

	rc, err := dcs.Download(f.OriginalID)
	if err != nil {
		return nil, err
	}

	t, err := tika.ExtractPlainContent(rc)
	if err != nil {
		return nil, err
	}

	f.Content = t.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ContentTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ContentTimestamp = &n
	}

	return f, nil
}

// processAudio uploads audio file to GCS and runs async speech to text over it
func (dfs *DropboxFs) processAudio(f *file.KazoupDropboxFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	bcs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.Dropbox)
	if err != nil {
		return nil, err
	}

	rc, err := bcs.Download(f.OriginalID)
	if err != nil {
		return nil, err
	}

	if err := objectstorage.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

	stt, err := sttlib.AsyncContent(fmt.Sprintf("gs://%s/%s", globals.AUDIO_BUCKET, f.ID))
	if err != nil {
		return nil, err
	}

	f.Content = stt.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			AudioTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.AudioTimestamp = &n
	}

	return f, nil
}

// processThumbnail, thumbnail generation
func (dfs *DropboxFs) processThumbnail(f *file.KazoupDropboxFile) (file.File, error) {
	// Downloads from dropbox, see connector
	dcs, err := cs.NewCloudStorageFromEndpoint(dfs.Endpoint, globals.Dropbox)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err = backoff.Retry(func() error {
		rc, err = dcs.Download(f.OriginalID)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	backoff.Retry(func() error {
		b, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			// Skip retry
			return nil
		}

		if err := objectstorage.Upload(ioutil.NopCloser(b), dfs.Endpoint.Index, f.ID); err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ThumbnailTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ThumbnailTimestamp = &n
	}

	return f, nil
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
		"file":"` + f.OriginalID + `",
		"include_inherited": true,
		"limit": 250
	}`)

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.DropboxFileMembers, bytes.NewBuffer(b))
	if err != nil {
		return f, err
	}
	req.Header.Set("Authorization", dfs.token())
	req.Header.Set("Content-Type", "application/json")
	rsp, err := c.Do(req)
	if err != nil {
		return f, err
	}
	defer rsp.Body.Close()

	var membersRsp *dropbox.FileMembersListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&membersRsp); err != nil {
		return f, err
	}

	if len(membersRsp.Users) > 0 {
		f.DropboxUsers = make([]dropbox.DropboxUser, 0)

		for _, v := range membersRsp.Users {
			a, err := dfs.getAccount(v.User.AccountID)
			if err != nil {
				return f, err
			}

			f.DropboxUsers = append(f.DropboxUsers, *a)
		}
	}

	if len(membersRsp.Invitees) > 0 {
		f.DropboxInvitees = membersRsp.Invitees
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
			if v.HasExplicitSharedMembers {
				f.Access = globals.ACCESS_SHARED

				// Error will be send over the channel withing the file
				f, err = dfs.getFileMembers(f)
			} else {
				// File is not share, but that means to dropbox that can be private, or public (everyone with link can access the file)
				for k, v := range dfs.PublicFiles {
					// Found
					if f.OriginalID == v.ID {
						f.Access = globals.ACCESS_PUBLIC

						// Remove found for performance and break
						dfs.PublicFiles = append(dfs.PublicFiles[:k], dfs.PublicFiles[k+1:]...)
						break
					}
				}
			}

			dfs.FilesChan <- NewFileMsg(f, err)
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
