package slack

import (
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

// Slack crawler
type Slack struct {
	Id       int64
	Running  chan bool
	Endpoint *proto_datasource.Endpoint
	Scanner  scan.Scanner
}

func NewSlack(id int64, dataSource *proto_datasource.Endpoint) *Slack {
	return &Slack{
		Id:       id,
		Running:  make(chan bool, 1),
		Endpoint: dataSource,
	}
}

// Start slack crawler
func (s *Slack) Start(crawls map[int64]scan.Scanner, ds int64) {
	log.Println("STARTTTTTT", s)
	go func() {

		// Slack scan finished
		s.Stop()
		delete(crawls, ds)
		s.sendCrawlerFinishedMsg()
	}()
}

// Stop slack crawler
func (s *Slack) Stop() {
	s.Running <- false
}

// Info slack crawler
func (s *Slack) Info() (scan.Info, error) {
	return scan.Info{
		Id:          s.Id,
		Type:        globals.Slack,
		Description: "Slack scanner",
	}, nil
}

func (s *Slack) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: s.Endpoint.Index,
	}

	if err := client.Publish(context.Background(), client.NewPublication(globals.CrawlerFinishedTopic, msg)); err != nil {
		return err
	}

	return nil
}
