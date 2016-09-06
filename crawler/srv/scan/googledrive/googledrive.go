package googledrive

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	googledrive "github.com/kazoup/platform/structs/googledrive"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"log"
	"time"
)

// Google crawler
type GoogleDrive struct {
	Id       int64
	Running  chan bool
	Endpoint *proto_datasource.Endpoint
	Scanner  scan.Scanner
}

func NewGoogleDrive(id int64, dataSource *proto_datasource.Endpoint) *GoogleDrive {
	return &GoogleDrive{
		Id:       id,
		Running:  make(chan bool, 1),
		Endpoint: dataSource,
	}
}

// Start google drive crawler
func (g *GoogleDrive) Start(crawls map[int64]scan.Scanner, ds int64) {
	go func() {
		if err := g.getFiles(); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		g.Stop()
		delete(crawls, ds)
		g.sendCrawlerFinishedMsg()
	}()
}

// Stop google drive crawler
func (g *GoogleDrive) Stop() {
	g.Running <- false
}

// Info google drive crawler
func (g *GoogleDrive) Info() (scan.Info, error) {
	return scan.Info{
		Id:          g.Id,
		Type:        globals.GoogleDrive,
		Description: "Google drive crawler",
	}, nil
}

func (s *GoogleDrive) getFiles() error {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  s.Endpoint.Token.AccessToken,
		TokenType:    s.Endpoint.Token.TokenType,
		RefreshToken: s.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(s.Endpoint.Token.Expiry, 0),
	})

	srv, err := drive.New(c)
	if err != nil {
		return err
	}

	r, err := srv.Files.List().PageSize(100).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		return sendFileMessagesForPage(r.Files, s.Endpoint.Index)
	}

	if len(r.NextPageToken) > 0 {
		return s.getNextPage(srv, r.NextPageToken)
	}

	return nil
}

func (s *GoogleDrive) getNextPage(srv *drive.Service, nextPageToken string) error {
	r, err := srv.Files.List().PageToken(nextPageToken).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		return sendFileMessagesForPage(r.Files, s.Endpoint.Index)
	}

	if len(r.NextPageToken) > 0 {
		return s.getNextPage(srv, r.NextPageToken)
	}

	return nil
}

func (g *GoogleDrive) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: g.Endpoint.Index,
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
		return err
	}

	return nil
}

func sendFileMessagesForPage(files []*drive.File, index string) error {
	for _, v := range files {
		f := googledrive.NewKazoupFileFromGoogleDriveFile(v)

		b, err := json.Marshal(f)
		if err != nil {
			return nil
		}

		msg := &crawler.FileMessage{
			Id:    getMD5Hash(v.WebViewLink),
			Index: index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.FilesTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
