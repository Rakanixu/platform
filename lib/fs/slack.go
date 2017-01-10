package fs

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/slack"
	"github.com/micro/go-micro/client"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// SlackFs slack file system
type SlackFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

// NewSlackFsFromEndpoint constructor
func NewSlackFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &SlackFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover files in slack datasource
func (sfs *SlackFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := sfs.getFiles(1); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		sfs.WalkRunning <- false
	}()

	return sfs.FilesChan, sfs.WalkRunning
}

// WalUsers discover users in slack
func (sfs *SlackFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		sfs.getUsers()

		// Slack user scan finished
		sfs.WalkUsersRunning <- false
	}()

	return sfs.UsersChan, sfs.WalkUsersRunning
}

// WalChannels discover channels in slack
func (sfs *SlackFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		sfs.getChannels()

		// Slack channels scan finished
		sfs.WalkChannelsRunning <- false
	}()

	return sfs.ChannelsChan, sfs.WalkChannelsRunning
}

// Token returns slack user token
func (sfs *SlackFs) Token(c client.Client) string {
	return "Bearer " + sfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (sfs *SlackFs) GetDatasourceId() string {
	return sfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (sfs *SlackFs) GetThumbnail(id string, c client.Client) (string, error) {
	return "", nil
}

// CreateFile belongs to Fs interface
func (sfs *SlackFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	return sfs.FilesChan
}

// DeleteFile deletes a slack file
func (sfs *SlackFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	return sfs.FilesChan
}

// ShareFile sets a PermalinkPublic available, so everyone with URL has access to the slack file
func (sfs *SlackFs) Update(req file_proto.ShareRequest) chan FileMsg {
	/*	if req.SharePublicly {
			return sfs.shareFilePublicly(req.OriginalId)
		} else {
			r := c.NewRequest(
				globals.DB_SERVICE_NAME,
				"DB.Read",
				&db_proto.ReadRequest{
					Index: req.Index,
					Type:  "file",
					Id:    req.FileId,
				},
			)
			rsp := &db_proto.ReadResponse{}
			if err := c.Call(ctx, r, rsp); err != nil {
				return "", err
			}

			var f *file.KazoupSlackFile
			if err := json.Unmarshal([]byte(rsp.Result), &f); err != nil {
				return "", err
			}

			return sfs.shareFileInsideTeam(f, req.DestinationId)
		}*/
	return sfs.FilesChan
}

// DownloadFile retrieves a file
func (sfs *SlackFs) DownloadFile(url string, opts ...string) (io.ReadCloser, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", sfs.token())
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// UploadFile uploads a file into google cloud storage
func (sfs *SlackFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, sfs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (sfs *SlackFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(sfs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (sfs *SlackFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(sfs.Endpoint.Index, "")
}

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
		pr, err := sfs.DownloadFile(sf.URLPrivateDownload)
		if err != nil {
			return errors.New("ERROR downloading slack file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for slack file")
		}

		if err := sfs.UploadFile(b, id); err != nil {
			return errors.New("ERROR uploading thumbnail for slack file")
		}
	}

	return nil
}

// shareFilePublicly will set a PermalinkPublic available and reachable for a file arcived/ stored in slack
func (sfs *SlackFs) shareFilePublicly(c client.Client, id string) (string, error) {
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
	f := file.NewKazoupFileFromSlackFile(ssr.File, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)
	if err := file.IndexAsync(c, f, globals.FilesTopic, sfs.Endpoint.Index, true); err != nil {
		return "", err
	}

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
