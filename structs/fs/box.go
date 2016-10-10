package fs

import (
	//"bytes"
	"encoding/json"
	//"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/box"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"log"
	"golang.org/x/oauth2"
	"net/http"
	//"net/url"
	"time"
)

type BoxFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
	Directories chan string
	LastDirTime int64
}

func NewBoxFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &BoxFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
		Directories: make(chan string),
	}
}

func (bfs *BoxFs) List() (chan file.File, chan bool, error) {
	bfs.refreshToken()

	go func() {
		bfs.LastDirTime = time.Now().Unix()
		for {
			select {
			case v := <-bfs.Directories:
				bfs.LastDirTime = time.Now().Unix()

				err := bfs.getDirChildren(v)
				if err != nil {
					log.Println(err)
				}
			default:
				// Helper for close channel and set that scanner has finish
				if bfs.LastDirTime+10 < time.Now().Unix() {
					bfs.Running <- false
					close(bfs.Directories)
					return
				}
			}

		}
	}()

	go func() {
		if err := bfs.getContentFromRoot(); err != nil {
			log.Println(err)
		}
	}()

	return bfs.FilesChan, bfs.Running, nil
}

func (bfs *BoxFs) Token() string {
	return "Bearer " + bfs.Endpoint.Token.AccessToken
}

func (bfs *BoxFs) GetDatasourceId() string {
	return bfs.Endpoint.Id
}

func (bfs *BoxFs) GetThumbnail(id string) (string, error) {
	/*args := `{"path":"` + id + `","size":{".tag":"w640h480"}}`
	url := fmt.Sprintf("%s?authorization=%s&arg=%s", globals.DropboxThumbnailEndpoint, dfs.Token(), url.QueryEscape(args))
*/
	return "", nil
}

func (bfs *BoxFs) getContentFromRoot() error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxFoldersEndpoint + "0", nil)
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

	log.Println("FILES")
	log.Println(bdc)


	for _, v := range bdc.ItemCollection.Entries {
		f := file.NewKazoupFileFromBoxFile(&v, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)
		bfs.FilesChan <- f

		if v.Type == "folder" {
			bfs.Directories <- v.ID
		}
	}

	return nil
}

// getDirChildren get children from directory
func (bfs *BoxFs) getDirChildren(id string) error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", globals.BoxFoldersEndpoint + id, nil)
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

	log.Println("FILES")
	log.Println(bdc)


	for _, v := range bdc.ItemCollection.Entries {
		f := file.NewKazoupFileFromBoxFile(&v, bfs.Endpoint.Id, bfs.Endpoint.UserId, bfs.Endpoint.Index)
		bfs.FilesChan <- f

		if v.Type == "folder" {
			bfs.Directories <- v.ID
		}
	}

	return nil
}

// refreshToken gets a new token from custom one and saves it
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

