package fs

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/osext"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	Drive  = "drive/"
	Drives = "drives/"
)

type OneDriveFs struct {
	Endpoint    *datasource_proto.Endpoint
	Running     chan bool
	FilesChan   chan file.File
	DrivesId    []string
	Directories chan string
	LastDirTime int64
}

func NewOneDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &OneDriveFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
		DrivesId:  []string{},
		// This is important to have a size bigger than one, the bigger, less likely to block
		// If not, program execution will block, due to recursivity,
		// We are pushing more elements before finish execution.
		// I expect to never push 10000 folders before other folders have been completly scanned
		Directories: make(chan string, 10000),
	}
}

func (ofs *OneDriveFs) List() (chan file.File, chan bool, error) {
	go func() {
		ofs.LastDirTime = time.Now().Unix()
		for {
			select {
			case v := <-ofs.Directories:
				ofs.LastDirTime = time.Now().Unix()

				err := ofs.getDirChildren(v)
				if err != nil {
					log.Println(err)
				}
			default:
				// Helper for close channel and set that scanner has finish
				if ofs.LastDirTime+10 < time.Now().Unix() {
					ofs.Running <- false
					close(ofs.Directories)
					return
				}
			}

		}
	}()

	go func() {
		if err := ofs.getFiles(); err != nil {
			log.Println(err)
		}
	}()

	return ofs.FilesChan, ofs.Running, nil
}

func (ofs *OneDriveFs) Token() string {
	return ofs.Endpoint.Token.AccessToken
}

func (ofs *OneDriveFs) GetDatasourceId() string {
	return ofs.Endpoint.Id
}

func (ofs *OneDriveFs) GetThumbnail(id string) (string, error) {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := fmt.Sprintf("%sitems/%s/thumbnails/0/medium", globals.OneDriveEndpoint+Drive, id)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var thumbRsp *onedrive.FileThumbnailResponse
	if err := json.NewDecoder(res.Body).Decode(&thumbRsp); err != nil {
		return "", err
	}

	return thumbRsp.URL, nil
}

func (ofs *OneDriveFs) CreateFile(fileType string) (string, error) {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

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
	// https://dev.onedrive.com/items/upload_put.htm
	url := fmt.Sprintf("%sroot:/untitled.%s:/content", globals.OneDriveEndpoint+Drive, globals.GetDocumentTemplate(fileType, false))
	req, err := http.NewRequest("PUT", url, t) // We require a template to be able to open / edit this files online
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	req.Header.Set("Content-Type", globals.ONEDRIVE_TEXT)
	if err != nil {
		return "", err
	}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var f *onedrive.OneDriveFile
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		return "", err
	}

	kfo := file.NewKazoupFileFromOneDriveFile(f, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)
	if err := file.IndexAsync(kfo, globals.FilesTopic, ofs.Endpoint.Index); err != nil {
		return "", err
	}

	return kfo.GetURL(), nil
}

// getFiles retrieves drives, directories and files
func (ofs *OneDriveFs) getFiles() error {
	if err := ofs.refreshToken(); err != nil {
		log.Println(err)
	}

	if err := ofs.getDrives(); err != nil {
		return err
	}
	if err := ofs.getDrivesChildren(); err != nil {
		return err
	}

	return nil
}

// refreshToken gets a new token from custom one and saves it
func (ofs *OneDriveFs) refreshToken() error {
	tokenSource := globals.NewMicrosoftOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  ofs.Endpoint.Token.AccessToken,
		TokenType:    ofs.Endpoint.Token.TokenType,
		RefreshToken: ofs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(ofs.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return err
	}
	ofs.Endpoint.Token.AccessToken = t.AccessToken
	ofs.Endpoint.Token.TokenType = t.TokenType
	ofs.Endpoint.Token.RefreshToken = t.RefreshToken
	ofs.Endpoint.Token.Expiry = t.Expiry.Unix()

	b, err := json.Marshal(ofs.Endpoint)
	if err != nil {
		return err
	}

	c := db_proto.NewDBClient("", nil)
	_, err = c.Update(context.Background(), &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    ofs.Endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}

// getDrives retrieve user drives
func (ofs *OneDriveFs) getDrives() error {
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := globals.OneDriveEndpoint + Drives
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var drivesRsp *onedrive.DrivesListResponse
	if err := json.NewDecoder(res.Body).Decode(&drivesRsp); err != nil {
		return err
	}

	for _, v := range drivesRsp.Value {
		ofs.DrivesId = append(ofs.DrivesId, v.ID)
	}

	return nil
}

// getDrivesChildren gets first level element from every found  drive
func (ofs *OneDriveFs) getDrivesChildren() error {
	var url string
	c := &http.Client{}

	for _, v := range ofs.DrivesId {
		//https://api.onedrive.com/v1.0/drives/f5a34c5d0f17415a/root/children
		url = globals.OneDriveEndpoint + Drives + v + "/root/children"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
		if err != nil {
			return err
		}
		res, err := c.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		var filesRsp *onedrive.FilesListResponse
		if err := json.NewDecoder(res.Body).Decode(&filesRsp); err != nil {
			return err
		}

		for _, v := range filesRsp.Value {
			// Is directory
			if len(v.File.MimeType) == 0 {
				ofs.Directories <- v.ID
				// Is file
			} else {
				f := file.NewKazoupFileFromOneDriveFile(&v, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)
				ofs.FilesChan <- f
			}
		}
	}

	return nil
}

// getDirChildren get children from directory
func (ofs *OneDriveFs) getDirChildren(id string) error {
	// https://api.onedrive.com/v1.0/drive/items/F5A34C5D0F17415A!114/children
	c := &http.Client{}
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/children"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", ofs.Endpoint.Token.TokenType+" "+ofs.Endpoint.Token.AccessToken)
	if err != nil {
		return err
	}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var filesRsp *onedrive.FilesListResponse
	if err := json.NewDecoder(res.Body).Decode(&filesRsp); err != nil {
		return err
	}

	for _, v := range filesRsp.Value {
		if len(v.File.MimeType) == 0 {
			//ofs.DirCounter++
			ofs.Directories <- v.ID
		} else {
			f := file.NewKazoupFileFromOneDriveFile(&v, ofs.Endpoint.Id, ofs.Endpoint.UserId, ofs.Endpoint.Index)
			ofs.FilesChan <- f
		}
	}

	return nil
}
