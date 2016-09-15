package onedrive

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	proto_db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	Drive  = "drive/"
	Drives = "drives/"
)

// OneDrive crawler
type OneDrive struct {
	Id           int64
	Running      chan bool
	Endpoint     *proto_datasource.Endpoint
	DrivesId     []string
	Direcotories chan string
	scan.Scanner
}

// NewOneDrive constructor
func NewOneDrive(id int64, dataSource *proto_datasource.Endpoint) *OneDrive {
	return &OneDrive{
		Id:           id,
		Running:      make(chan bool, 1),
		Endpoint:     dataSource,
		DrivesId:     []string{},
		Direcotories: make(chan string, 100),
	}
}

// Start scan
func (o *OneDrive) Start(crawls map[int64]scan.Scanner, ds int64) {
	go func() {
		if err := o.getFiles(); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 5)
		if err := o.clearIndex(); err != nil {
			log.Println(err)
		}

		// One drive scan finished
		o.Stop()
		delete(crawls, ds)
		o.sendCrawlerFinishedMsg()
	}()
}

// Stop scan
func (o *OneDrive) Stop() {
	o.Running <- false
}

// Info about scan
func (o *OneDrive) Info() (scan.Info, error) {
	return scan.Info{
		Id:          o.Id,
		Type:        globals.OneDrive,
		Description: "One drive crawler",
	}, nil
}

// Compares LastSeen with the time the crawler started
// so all records with a LastSeen before will be removed from index
// file does not exists any more on datasource
func (o *OneDrive) clearIndex() error {
	c := proto_db.NewDBClient("", nil)
	_, err := c.DeleteByQuery(context.Background(), &proto_db.DeleteByQueryRequest{
		Indexes:  []string{o.Endpoint.Index},
		Types:    []string{"file"},
		LastSeen: o.Endpoint.LastScanStarted,
	})
	if err != nil {
		return err
	}

	return nil
}

// getFiles retrieves drives, directories and files
func (o *OneDrive) getFiles() error {
	if err := o.refreshToken(); err != nil {
		log.Println(err)
	}

	o.listenForDirs()
	if err := o.getDrives(); err != nil {
		return err
	}
	if err := o.getDrivesChildren(); err != nil {
		return err
	}

	return nil
}

// getDrives retrieve user drives
func (o *OneDrive) getDrives() error {
	c := &http.Client{}
	//https://api.onedrive.com/v1.0/drives
	url := globals.OneDriveEndpoint + Drives
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", o.Endpoint.Token.TokenType+" "+o.Endpoint.Token.AccessToken)
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
		o.DrivesId = append(o.DrivesId, v.ID)
	}

	return nil
}

// getDrivesChildren gets first level element from every found  drive
func (o *OneDrive) getDrivesChildren() error {
	var url string
	c := &http.Client{}

	for _, v := range o.DrivesId {
		//https://api.onedrive.com/v1.0/drives/f5a34c5d0f17415a/root/children
		url = globals.OneDriveEndpoint + Drives + v + "/root/children"

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", o.Endpoint.Token.TokenType+" "+o.Endpoint.Token.AccessToken)
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
				o.Direcotories <- v.ID
				// Is file
			} else {
				f := file.NewKazoupFileFromOneDriveFile(&v)
				err := o.sendFileMsg(*f, v.WebURL)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

// listenForDirs listen for directories id to retrieve its contents
func (o *OneDrive) listenForDirs() error {
	go func() {
		for {
			select {
			case v := <-o.Direcotories:
				err := o.getDirChildren(v)
				if err != nil {
					log.Println(err)
				}
			}

		}
	}()

	return nil
}

// getDirChildren get children from directory
func (o *OneDrive) getDirChildren(id string) error {
	// https://api.onedrive.com/v1.0/drive/items/F5A34C5D0F17415A!114/children
	c := &http.Client{}
	url := globals.OneDriveEndpoint + Drive + "items/" + id + "/children"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", o.Endpoint.Token.TokenType+" "+o.Endpoint.Token.AccessToken)
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
			o.Direcotories <- v.ID
		} else {
			f := file.NewKazoupFileFromOneDriveFile(&v)
			err := o.sendFileMsg(*f, v.WebURL)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// sendFileMsg publishes file messages
func (o *OneDrive) sendFileMsg(f interface{}, url string) error {
	b, err := json.Marshal(f)
	if err != nil {
		return nil
	}

	msg := &crawler.FileMessage{
		Id:    getMD5Hash(url),
		Index: o.Endpoint.Index,
		Data:  string(b),
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.FilesTopic, msg)); err != nil {
		return err
	}
	return nil
}

// sendCrawlerFinishedMsg publishes crawler finished messages
func (o *OneDrive) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: o.Endpoint.Id,
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
		return err
	}

	return nil
}

// refreshToken gets a new token from custom one and saves it
func (o *OneDrive) refreshToken() error {
	tokenSource := globals.NewMicrosoftOauthConfig().TokenSource(oauth2.NoContext, &oauth2.Token{
		AccessToken:  o.Endpoint.Token.AccessToken,
		TokenType:    o.Endpoint.Token.TokenType,
		RefreshToken: o.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(o.Endpoint.Token.Expiry, 0),
	})

	t, err := tokenSource.Token()
	if err != nil {
		return err
	}
	o.Endpoint.Token.AccessToken = t.AccessToken
	o.Endpoint.Token.TokenType = t.TokenType
	o.Endpoint.Token.RefreshToken = t.RefreshToken
	o.Endpoint.Token.Expiry = t.Expiry.Unix()

	b, err := json.Marshal(o.Endpoint)
	if err != nil {
		return err
	}

	c := proto_db.NewDBClient("", nil)
	_, err = c.Update(context.Background(), &proto_db.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    o.Endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
