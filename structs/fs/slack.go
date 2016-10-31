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

// SlackFs slack file system
type SlackFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

// NewSlackFsFromEndpoint constructor
func NewSlackFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &SlackFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

// List returns 2 channels, for files and state. Discover files in slack datasource
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

// Token returns slack user token
func (sfs *SlackFs) Token() string {
	return "Bearer " + sfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (sfs *SlackFs) GetDatasourceId() string {
	return sfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (sfs *SlackFs) GetThumbnail(id string) (string, error) {
	return "", nil
}

// CreateFile belongs to Fs interface
func (sfs *SlackFs) CreateFile(fileType string) (string, error) {
	return "", nil
}

// ShareFile sets a PermalinkPublic available, so everyone with URL has access to the slack file
func (sfs *SlackFs) ShareFile(id string, sharePublicly bool) (string, error) {
	if sharePublicly {
		return sfs.shareFilePublicly(id)
	} else {
		return sfs.shareFileInsideTeam(id)
	}
}

// getUsers retrieves users from slack team
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

// getChannels retrieves channels from slack team
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

// getFiles discover slack files for user and push them to broker
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
		f := file.NewKazoupFileFromSlackFile(&v, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)

		sfs.FilesChan <- f
	}

	if filesRsp.Paging.Pages >= page {
		sfs.getFiles(page + 1)
	}

	return nil
}

// shareFilePublicly will set a PermalinkPublic available and reachable for a file arcived/ stored in slack
func (sfs *SlackFs) shareFilePublicly(id string) (string, error) {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)
	data.Add("file", id)

	c := &http.Client{}
	rsp, err := c.PostForm(globals.SlackShareFilesEndpoint, data)

	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var ssr *slack.SlackShareResponse
	if err := json.NewDecoder(rsp.Body).Decode(&ssr); err != nil {
		return "", err
	}

	log.Println("SHARE PUBLIC")
	log.Println(rsp.StatusCode)

	// Response contains object, permalink_public attr will be modified
	// Reindex document
	f := file.NewKazoupFileFromSlackFile(&ssr.File, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)
	if err := file.IndexAsync(f, globals.FilesTopic, sfs.Endpoint.Index); err != nil {
		return "", err
	}

	return ssr.File.PermalinkPublic, nil
}

// shareFileInsideTeam will post a message to a channel or team member linking the slack file
func (sfs *SlackFs) shareFileInsideTeam(id string) (string, error) {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)
	data.Add("channel", "C02EW2K5U")
	data.Add("attachments", `[
		{
		    "fallback": "fallback attr",
		    "pretext": "Kazoup backend test, sharing / mentioning a file in slack with a slack channel #ramdom",
		    "title": "step.mov",
		    "title_link": "https://kazoup.slack.com/files/pablo.aguirre/F2WBHHNKB/step.mov",
		    "text": "Yeap, it is working mate! Long live Kazoup",
		    "color": "#7CD197"
		}
	]`)

	c := &http.Client{}
	rsp, err := c.PostForm(globals.SlackPostMessageEndpoint, data)

	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	log.Println("SHARE PRIVATE")
	log.Println(rsp.StatusCode)

	/*	var ssr *slack.SlackShareResponse
		if err := json.NewDecoder(rsp.Body).Decode(&ssr); err != nil {
			return "", err
		}

		// Response contains object, permalink_public attr will be modified
		// Reindex document
		f := file.NewKazoupFileFromSlackFile(&ssr.File, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)
		if err := file.IndexAsync(f, globals.FilesTopic, sfs.Endpoint.Index); err != nil {
			return "", err
		}

		return ssr.File.PermalinkPublic, nil*/
	return "", nil
}
