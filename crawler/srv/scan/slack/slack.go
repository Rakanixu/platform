package slack

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/crawler/srv/scan"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/slack"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
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
	go func() {
		if err := s.getUsers(); err != nil {
			log.Println(err)
		}

		if err := s.getChannels(); err != nil {
			log.Println(err)
		}

		if err := s.getFiles(1); err != nil {
			log.Println(err)
		}
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

func (s *Slack) getUsers() error {
	data := make(url.Values)
	data.Add("token", s.Endpoint.Token.AccessToken)

	c := &http.Client{}

	rsp, err := c.PostForm(globals.SlackUsersEndpoint, data)

	if err != nil {

		return err
	}
	defer rsp.Body.Close()

	var usersRsp *slack.UserListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&usersRsp); err != nil {
		return err
	}

	for _, v := range usersRsp.Members {
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}

		msg := &crawler.SlackUserMessage{
			Id:    v.ID,
			Index: s.Endpoint.Index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.SlackUsersTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Slack) getChannels() error {
	data := make(url.Values)
	data.Add("token", s.Endpoint.Token.AccessToken)

	c := &http.Client{}

	rsp, err := c.PostForm(globals.SlackChannelsEndpoint, data)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	var channelsRsp *slack.ChannelListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&channelsRsp); err != nil {
		return err
	}

	for _, v := range channelsRsp.Channels {
		b, err := json.Marshal(v)
		if err != nil {
			return nil
		}

		msg := &crawler.SlackChannelMessage{
			Id:    v.ID,
			Index: s.Endpoint.Index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.SlackChannelsTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Slack) getFiles(page int) error {
	data := make(url.Values)
	data.Add("token", s.Endpoint.Token.AccessToken)
	data.Add("page", strconv.Itoa(page))

	c := &http.Client{}

	rsp, err := c.PostForm(globals.SlackFilesEndpoint, data)
	if err != nil {

		return err
	}
	defer rsp.Body.Close()

	var filesRsp *slack.FilesListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&filesRsp); err != nil {
		return err
	}

	for _, v := range filesRsp.Files {
		f := file.NewKazoupFileFromSlackFile(&v, s.Endpoint.Id)

		b, err := json.Marshal(f)
		if err != nil {
			log.Println(err)
		}

		msg := &crawler.FileMessage{
			Id:    getMD5Hash(v.URLPrivate),
			Index: s.Endpoint.Index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.FilesTopic, msg)); err != nil {
			return err
		}
	}

	if filesRsp.Paging.Pages >= page {
		s.getFiles(page + 1)
	}

	return nil
}

func (s *Slack) sendCrawlerFinishedMsg() error {
	msg := &crawler.CrawlerFinishedMessage{
		DatasourceId: s.Endpoint.Id,
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
