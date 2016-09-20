package googledrive

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

// Google crawler
type GoogleDrive struct {
	Id       int64
	Running  chan bool
	Endpoint *datasource_proto.Endpoint
	Scanner  scan.Scanner
}

func NewGoogleDrive(id int64, dataSource *datasource_proto.Endpoint) *GoogleDrive {
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
		time.Sleep(time.Second * 5)
		if err := g.clearIndex(); err != nil {
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

// Compares LastSeen with the time the crawler started
// so all records with a LastSeen before will be removed from index
// file does not exists any more on datasource
func (g *GoogleDrive) clearIndex() error {
	c := db_proto.NewDBClient("", nil)
	_, err := c.DeleteByQuery(context.Background(), &db_proto.DeleteByQueryRequest{
		Indexes:  []string{g.Endpoint.Index},
		Types:    []string{"file"},
		LastSeen: g.Endpoint.LastScanStarted,
	})
	if err != nil {
		return err
	}

	return nil
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
		if err := sendFileMessagesForPage(r.Files, s.Endpoint); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := s.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

func (s *GoogleDrive) getNextPage(srv *drive.Service, nextPageToken string) error {
	r, err := srv.Files.List().PageToken(nextPageToken).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := sendFileMessagesForPage(r.Files, s.Endpoint); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := s.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoogleDrive) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: g.Endpoint.Id,
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
		return err
	}

	return nil
}

func sendFileMessagesForPage(files []*drive.File, ds *datasource_proto.Endpoint) error {
	for _, v := range files {
		//Conflicts with ES and size in  google is parse as string
		f := file.NewKazoupFileFromGoogleDriveFile(v, ds.Id)

		b, err := json.Marshal(f)
		if err != nil {
			return nil
		}
		msg := &crawler.FileMessage{
			Id:    getMD5Hash(v.WebViewLink),
			Index: ds.Index,
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
