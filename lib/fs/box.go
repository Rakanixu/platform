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
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/micro/go-micro/client"
	"golang.org/x/oauth2"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

// BoxFs Box File System
type BoxFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan file.File
	FileMetaChan        chan FileMsg
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
		FilesChan:           make(chan file.File),
		FileMetaChan:        make(chan FileMsg),
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
func (bfs *BoxFs) Walk() (chan file.File, chan bool, error) {
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

	return bfs.FilesChan, bfs.WalkRunning, nil
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

// Token returns user token for box datasource
func (bfs *BoxFs) Token(c client.Client) string {
	bfs.refreshToken(c)

	return "Bearer " + bfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (bfs *BoxFs) GetDatasourceId() string {
	return bfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to an image
func (bfs *BoxFs) GetThumbnail(id string, c client.Client) (string, error) {
	url := fmt.Sprintf(
		"%s%s&Authorization=%s",
		globals.BoxFileMetadataEndpoint,
		id,
		"/thumbnail.png?min_height=256&min_width=256",
		bfs.Token(c),
	)

	return url, nil
}

// Create file in box
func (bfs *BoxFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		// Box supports multi part form upload
		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		// File template path
		t := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(rq.MimeType, true))
		buf := &bytes.Buffer{}
		mw := multipart.NewWriter(buf)
		defer mw.Close()

		f, err := os.Open(t)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		defer f.Close()

		// This is how you upload a file as multipart form
		ff, err := mw.CreateFormFile("file", t)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		if _, err = io.Copy(ff, f); err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		// Add extra fields required by API
		mw.WriteField(
			"attributes",
			`{"name":"`+rq.FileName+`.`+globals.GetDocumentTemplate(rq.MimeType, false)+`", "parent":{"id":"0"}}`,
		)
		if err := mw.Close(); err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		hc := &http.Client{}
		req, err := http.NewRequest("POST", globals.BoxUploadEndpoint, buf)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		req.Header.Set("Authorization", bfs.token())
		req.Header.Set("Content-Type", mw.FormDataContentType())

		rsp, err := hc.Do(req)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		var bf *box.BoxUpload
		if err := json.NewDecoder(rsp.Body).Decode(&bf); err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		if rsp.StatusCode == http.StatusConflict {
			bfs.FileMetaChan <- NewFileMsg(nil, errors.New("Conflict creating file in Box, file with same name already exists"))
			return
		}

		if rsp.StatusCode != http.StatusCreated && bf.TotalCount != 1 {
			bfs.FileMetaChan <- NewFileMsg(nil, errors.New("Failed creating file in Box"))
			return
		}

		// Construct Kazoup file from box created file and index it
		kfb := file.NewKazoupFileFromBoxFile(bf.Entries[0], bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

		bfs.FileMetaChan <- NewFileMsg(kfb, nil)
	}()

	return bfs.FileMetaChan
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
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", bfs.token())
		rsp, err := bc.Do(r)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusNoContent {
			bfs.FileMetaChan <- NewFileMsg(nil, errors.New(fmt.Sprintf("Deleting Box file failed with status code %d", rsp.StatusCode)))
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		bfs.FileMetaChan <- NewFileMsg(
			&file.KazoupBoxFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return bfs.FileMetaChan
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
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		r.Header.Set("Authorization", bfs.token())
		rsp, err := bc.Do(r)
		if err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}
		defer rsp.Body.Close()

		if rsp.StatusCode != http.StatusOK {
			bfs.FileMetaChan <- NewFileMsg(nil, errors.New(fmt.Sprintf("Sharing Box file failed with status code %d", rsp.StatusCode)))
			return
		}

		var f *box.BoxFileMeta
		if err := json.NewDecoder(rsp.Body).Decode(&f); err != nil {
			bfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		kbf := file.NewKazoupFileFromBoxFile(*f, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

		bfs.FileMetaChan <- NewFileMsg(kbf, nil)
	}()

	return bfs.FileMetaChan
}

// DownloadFile retrieves a file
func (bfs *BoxFs) DownloadFile(id string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/content", globals.BoxFileMetadataEndpoint, id)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(r)
	if err != nil {
		return nil, err
	}

	return rsp.Body, nil
}

// UploadFile uploads a file into google cloud storage
func (bfs *BoxFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, bfs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (bfs *BoxFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(bfs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (bfs *BoxFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(bfs.Endpoint.Index, "")
}

// getDirChildren get children from directory
func (bfs *BoxFs) getDirChildren(id string, offset, limit int) error {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/?offset=%d&limit=%d", globals.BoxFoldersEndpoint, id, offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var bdc *box.BoxDirContents
	if err := json.NewDecoder(rsp.Body).Decode(&bdc); err != nil {
		return err
	}

	for _, v := range bdc.ItemCollection.Entries {
		if v.Type == "folder" {
			// Push found directories into the queue to be crawled
			bfs.Directories <- v.ID
		} else {
			// File discovered, but need to retrieve more info about the file
			if err := bfs.getMetadataFromFile(v.ID); err != nil {
				return err
			}
		}
	}

	if bdc.ItemCollection.TotalCount > bdc.ItemCollection.Offset+bdc.ItemCollection.Limit {
		bfs.getDirChildren(
			id,
			bdc.ItemCollection.Offset+bdc.ItemCollection.Limit,
			bdc.ItemCollection.Limit,
		)
	}

	return nil
}

// getMetadataFromFile retrieves more info about discovered files in box
func (bfs *BoxFs) getMetadataFromFile(id string) error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxFileMetadataEndpoint+id, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var fm *box.BoxFileMeta
	if err := json.NewDecoder(rsp.Body).Decode(&fm); err != nil {
		return err
	}

	f := file.NewKazoupFileFromBoxFile(*fm, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)

	if err := bfs.generateThumbnail(fm, f.ID); err != nil {
		log.Println(err)
	}

	bfs.FilesChan <- f

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage with kazoup id
func (bfs *BoxFs) generateThumbnail(fm *box.BoxFileMeta, id string) error {
	name := strings.Split(fm.Name, ".")

	if categories.GetDocType("."+name[len(name)-1]) == globals.CATEGORY_PICTURE {
		rc, err := bfs.DownloadFile(fm.ID)
		if err != nil {
			return errors.New("ERROR downloading box file")
		}

		rd, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for box file")
		}

		if err := bfs.UploadFile(rd, id); err != nil {
			return errors.New("ERROR uploading thumbnail for box file")
		}
	}

	return nil
}

// refreshToken gets a new token (refreshed if expired) from custom one and saves it
func (bfs *BoxFs) refreshToken(cl client.Client) error {
	tokenSource := globals.NewBoxOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  bfs.Endpoint.Token.AccessToken,
		TokenType:    bfs.Endpoint.Token.TokenType,
		RefreshToken: bfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(bfs.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return err
	}
	bfs.Endpoint.Token.AccessToken = t.AccessToken
	bfs.Endpoint.Token.TokenType = t.TokenType
	bfs.Endpoint.Token.RefreshToken = t.RefreshToken
	bfs.Endpoint.Token.Expiry = t.Expiry.Unix()

	b, err := json.Marshal(bfs.Endpoint)
	if err != nil {
		return err
	}

	req := cl.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Update",
		&db_proto.UpdateRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    bfs.Endpoint.Id,
			Data:  string(b),
		},
	)
	rsp := &db_proto.UpdateResponse{}
	if err := cl.Call(globals.NewSystemContext(), req, rsp); err != nil {
		return err
	}

	return nil
}

// token returns authorization header token as string
func (bfs *BoxFs) token() string {
	return "Bearer " + bfs.Endpoint.Token.AccessToken
}
