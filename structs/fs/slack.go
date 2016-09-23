package fs

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/slack"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type SlackFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

func NewSlackFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &SlackFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

func (sfs *SlackFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := sfs.getUsers(); err != nil {
			log.Println(err)
		}

		if err := sfs.getChannels(); err != nil {
			log.Println(err)
		}

		if err := sfs.getFiles(1); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		sfs.Running <- false
	}()

	return sfs.FilesChan, sfs.Running, nil
}

func (sfs *SlackFs) Token() string {
	return sfs.Endpoint.Token.AccessToken
}

func (sfs *SlackFs) GetDatasourceId() string {
	return sfs.Endpoint.Id
}

func (sfs *SlackFs) getUsers() error {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)

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
			Index: sfs.Endpoint.Index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.SlackUsersTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func (sfs *SlackFs) getChannels() error {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)

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
			Index: sfs.Endpoint.Index,
			Data:  string(b),
		}

		if err := client.Publish(context.Background(), client.NewPublication(globals.SlackChannelsTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func (sfs *SlackFs) getFiles(page int) error {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)
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
		f := file.NewKazoupFileFromSlackFile(&v, sfs.Endpoint.Id, sfs.Endpoint.Index)

		sfs.FilesChan <- f
	}

	if filesRsp.Paging.Pages >= page {
		sfs.getFiles(page + 1)
	}

	return nil
}
