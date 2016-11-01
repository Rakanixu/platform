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
	"github.com/kazoup/platform/structs/box"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// BoxFs Box File System
type BoxFs struct {
	Endpoint      *datasource_proto.Endpoint
	Running       chan bool
	FilesChan     chan file.File
	Directories   chan string
	LastDirTime   int64
	DefaultOffset int
	DefaultLimit  int
}

// NewBoxFsFromEndpoint constructor
func NewBoxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &BoxFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
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
func (bfs *BoxFs) List() (chan file.File, chan bool, error) {
	bfs.refreshToken()

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
					bfs.Running <- false
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

	return bfs.FilesChan, bfs.Running, nil
}

// Token returns user token for box datasource
func (bfs *BoxFs) Token() string {
	bfs.refreshToken()

	return "Bearer " + bfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (bfs *BoxFs) GetDatasourceId() string {
	return bfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to an image
func (bfs *BoxFs) GetThumbnail(id string) (string, error) {
	url := fmt.Sprintf(
		"%s%s&Authorization=%s",
		globals.BoxFileMetadataEndpoint,
		id,
		"/thumbnail.png?min_height=256&min_width=256",
		bfs.Token(),
	)

	return url, nil
}

// CreateFile in box
func (bfs *BoxFs) CreateFile(fileType string) (string, error) {
	// Box supports multi part form upload
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		return "", err
	}

	// File template path
	t := fmt.Sprintf("%s%s%s", folderPath, "/doc_templates/", globals.GetDocumentTemplate(fileType, true))
	buf := &bytes.Buffer{}
	mw := multipart.NewWriter(buf)
	defer mw.Close()

	f, err := os.Open(t)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// This is how you upload a file as multipart form
	ff, err := mw.CreateFormFile("file", t)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(ff, f); err != nil {
		return "", err
	}

	// Add extra fields required by API
	mw.WriteField(
		"attributes",
		`{"name":"untitle.`+globals.GetDocumentTemplate(fileType, false)+`", "parent":{"id":"0"}}`,
	)
	if err := mw.Close(); err != nil {
		return "", err
	}

	c := &http.Client{}
	req, err := http.NewRequest("POST", globals.BoxUploadEndpoint, buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", bfs.Token())
	req.Header.Set("Content-Type", mw.FormDataContentType())

	rsp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var bf *box.BoxUpload
	if err := json.NewDecoder(rsp.Body).Decode(&bf); err != nil {
		return "", err
	}

	if rsp.StatusCode == http.StatusConflict {
		return "", errors.New("Conflict creating file in Box, file with same name already exists")
	}

	if rsp.StatusCode != http.StatusCreated && bf.TotalCount != 1 {
		return "", errors.New("Failed creating file in Box")
	}

	// Construct Kazoup file from box created file and index it
	kfb := file.NewKazoupFileFromBoxFile(&bf.Entries[0], bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)
	if err := file.IndexAsync(kfb, globals.FilesTopic, bfs.Endpoint.Index); err != nil {
		return "", err
	}

	return kfb.GetURL(), nil
}

// ShareFile
func (bfs *BoxFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	return "", nil
}

// getDirChildren get children from directory
func (bfs *BoxFs) getDirChildren(id string, offset, limit int) error {
	c := &http.Client{}
	url := fmt.Sprintf("%s%s/?offset=%d&limit=%d", globals.BoxFoldersEndpoint, id, offset, limit)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", bfs.Token())
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
	req.Header.Set("Authorization", bfs.Token())
	rsp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var fm *box.BoxFileMeta
	if err := json.NewDecoder(rsp.Body).Decode(&fm); err != nil {
		return err
	}

	f := file.NewKazoupFileFromBoxFile(fm, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)
	bfs.FilesChan <- f

	return nil
}

// refreshToken gets a new token (refreshed if expired) from custom one and saves it
func (bfs *BoxFs) refreshToken() error {
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

	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, nil)
	_, err = c.Update(globals.NewSystemContext(), &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    bfs.Endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}
