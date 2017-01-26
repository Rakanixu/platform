package fs

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/categories"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/slack"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// getUsers retrieves users from slack team
func (sfs *SlackFs) getUsers() {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)

	c := &http.Client{}

	rsp, err := c.PostForm(globals.SlackUsersEndpoint, data)

	if err != nil {
		sfs.UsersChan <- NewUserMsg(nil, err)
	}
	defer rsp.Body.Close()

	var usersRsp *slack.UserListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&usersRsp); err != nil {
		sfs.UsersChan <- NewUserMsg(nil, err)
	}

	for _, v := range usersRsp.Members {
		b, err := json.Marshal(v)
		if err != nil {
			sfs.UsersChan <- NewUserMsg(nil, err)
		}

		msg := &crawler.SlackUserMessage{
			Id:    v.ID,
			Index: sfs.Endpoint.Index,
			Data:  string(b),
		}

		sfs.UsersChan <- NewUserMsg(msg, nil)
	}
}

// getChannels retrieves channels from slack team
func (sfs *SlackFs) getChannels() {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)

	c := &http.Client{}

	rsp, err := c.PostForm(globals.SlackChannelsEndpoint, data)
	if err != nil {
		sfs.ChannelsChan <- NewChannelMsg(nil, err)
	}
	defer rsp.Body.Close()

	var channelsRsp *slack.ChannelListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&channelsRsp); err != nil {
		sfs.ChannelsChan <- NewChannelMsg(nil, err)
	}

	for _, v := range channelsRsp.Channels {
		b, err := json.Marshal(v)
		if err != nil {
			sfs.ChannelsChan <- NewChannelMsg(nil, err)
		}

		msg := &crawler.SlackChannelMessage{
			Id:    v.ID,
			Index: sfs.Endpoint.Index,
			Data:  string(b),
		}

		sfs.ChannelsChan <- NewChannelMsg(msg, nil)
	}
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
		f := file.NewKazoupFileFromSlackFile(v, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)

		if err := sfs.generateThumbnail(v, f.ID); err != nil {
			log.Println(err)
		}

		sfs.FilesChan <- NewFileMsg(f, nil)
	}

	if filesRsp.Paging.Pages >= page {
		sfs.getFiles(page + 1)
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (sfs *SlackFs) generateThumbnail(sf slack.SlackFile, id string) error {
	if categories.GetDocType("."+sf.Filetype) == globals.CATEGORY_PICTURE {
		// Download file from Slack, so connector is globals.Slack
		scs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.Slack)
		if err != nil {
			return err
		}
		pr, err := scs.Download(sf.URLPrivateDownload)
		if err != nil {
			return err
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return err
		}

		// Upload file to GoogleCloudStorage, so connector is globals.GoogleCloudStorage
		ncs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.GoogleCloudStorage)
		if err != nil {
			return err
		}

		if err := ncs.Upload(b, id); err != nil {
			return err
		}
	}

	return nil
}

// shareFilePublicly will set a PermalinkPublic available and reachable for a file arcived/ stored in slack
func (sfs *SlackFs) shareFilePublicly(id string) (string, error) {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)
	data.Add("file", id)

	hc := &http.Client{}
	rsp, err := hc.PostForm(globals.SlackShareFilesEndpoint, data)

	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var ssr *slack.SlackShareResponse
	if err := json.NewDecoder(rsp.Body).Decode(&ssr); err != nil {
		return "", err
	}

	// Response contains object, permalink_public attr will be modified
	// Reindex document
	//TODO: Commented out because client dependency. this method will return a FileMsg over a channel
	/*f := file.NewKazoupFileFromSlackFile(ssr.File, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)
	if err := file.IndexAsync(c, f, globals.FilesTopic, sfs.Endpoint.Index, true); err != nil {
		return "", err
	}*/

	return ssr.File.PermalinkPublic, nil
}

// shareFileInsideTeam will post a message to a channel or team member linking the slack file
func (sfs *SlackFs) shareFileInsideTeam(f *file.KazoupSlackFile, destId string) (string, error) {
	data := make(url.Values)
	data.Add("token", sfs.Endpoint.Token.AccessToken)
	data.Add("channel", destId)
	data.Add("attachments", `[
		{
		    "title": "`+f.Original.Name+`",
		    "title_link": "`+f.Original.Permalink+`",
		    "color": "#21a9f5"
		}
	]`)

	c := &http.Client{}
	rsp, err := c.PostForm(globals.SlackPostMessageEndpoint, data)

	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	return "", nil
}

func (sfs *SlackFs) token() string {
	return "Bearer " + sfs.Endpoint.Token.AccessToken
}
