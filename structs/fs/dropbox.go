package fs

import (
	"bytes"
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/dropbox"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"log"
	"net/http"
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
	/*	cfg := globals.NewGoogleOautConfig()
		c := cfg.Client(context.Background(), &oauth2.Token{
			AccessToken:  dfs.Endpoint.Token.AccessToken,
			TokenType:    dfs.Endpoint.Token.TokenType,
			RefreshToken: dfs.Endpoint.Token.RefreshToken,
			Expiry:       time.Unix(dfs.Endpoint.Token.Expiry, 0),
		})

		srv, err := drive.New(c)
		if err != nil {
			return "", err
		}
		r, err := srv.Files.Get(id).Fields("thumbnailLink").Do()
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%ss700", r.ThumbnailLink[:len(r.ThumbnailLink)-4]), nil*/
	return "", nil
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
		dfs.FilesChan <- f
	}

	return nil
}
