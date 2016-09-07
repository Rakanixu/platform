package onedrive

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

const (
	Drive  = "drive/"
	Drives = "drives/"
)

type OneDrive struct {
	Id           int64
	Running      chan bool
	Endpoint     *proto_datasource.Endpoint
	DrivesId     []string
	Direcotories chan string
	scan.Scanner
}

func NewOneDrive(id int64, dataSource *proto_datasource.Endpoint) *OneDrive {
	return &OneDrive{
		Id:           id,
		Running:      make(chan bool, 1),
		Endpoint:     dataSource,
		DrivesId:     []string{},
		Direcotories: make(chan string, 100),
	}
}

func (o *OneDrive) Start(crawls map[int64]scan.Scanner, ds int64) {
	go func() {
		if err := o.getFiles(); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		o.Stop()
		delete(crawls, ds)
		o.sendCrawlerFinishedMsg()
	}()
}

func (o *OneDrive) Stop() {
	o.Running <- false
}

func (o *OneDrive) Info() (scan.Info, error) {
	return scan.Info{
		Id:          o.Id,
		Type:        globals.OneDrive,
		Description: "One drive crawler",
	}, nil
}

func (o *OneDrive) getFiles() error {
	o.listenForDirs()

	o.getDrives()
	o.getDrivesChildren()

	return nil
}

func (o *OneDrive) getDrives() error {
	c := &http.Client{}
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
			f := onedrive.NewKazoupFileFromOneDriveFile(&v)
			err := o.sendFileMsg(*f, v.WebURL)
			if err != nil {
				return err
			}
			log.Println("??????", v.File.MimeType, len(v.File.MimeType))
			// Is directory
			if len(v.File.MimeType) == 0 {
				log.Println("IS DIR", v)
				o.Direcotories <- v.ID
			}
		}
	}

	return nil
}

func (o *OneDrive) listenForDirs() error {
	go func() {
		for {
			select {
			case v := <-o.Direcotories:
				log.Println("listenForDirs", v)
				o.getDirChildren(v)
			}

		}
	}()

	return nil
}

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
		f := onedrive.NewKazoupFileFromOneDriveFile(&v)
		err := o.sendFileMsg(*f, v.WebURL)
		if err != nil {
			return err
		}

		if len(v.File.MimeType) == 0 {
			o.Direcotories <- v.ID
		}
	}

	return nil
}

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

func (o *OneDrive) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: o.Endpoint.Index,
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
		return err
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
