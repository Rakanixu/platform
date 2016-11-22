package fs

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/slack"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"io/ioutil"
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
func (sfs *SlackFs) CreateFile(rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	return &file_proto.CreateResponse{}, nil
}

// DeleteFile deletes a slack file
func (sfs *SlackFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	return &file_proto.DeleteResponse{}, nil
}

// ShareFile sets a PermalinkPublic available, so everyone with URL has access to the slack file
func (sfs *SlackFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	if req.SharePublicly {
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
	}
}

// DownloadFile retrieves a file
func (sfs *SlackFs) DownloadFile(url string, opts ...string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", sfs.Token())
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// UploadFile uploads a file into google cloud storage
func (sfs *SlackFs) UploadFile(file []byte, fId string) error {
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
		if categories.GetDocType("."+v.Filetype) == globals.CATEGORY_PICTURE {
			b, err := sfs.DownloadFile(v.URLPrivateDownload)
			if err != nil {
				log.Println("ERROR downloading slack file: ", err)
			}

			b, err = image.Thumbnail(b, globals.THUMBNAIL_WIDTH)
			if err != nil {
				log.Println("ERROR generating thumbnail for slack file: ", err)
			}

			if err := sfs.UploadFile(b, v.ID); err != nil {
				log.Println("ERROR uploading thumbnail for slack file: ", err)
			}
		}

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

	// Response contains object, permalink_public attr will be modified
	// Reindex document
	f := file.NewKazoupFileFromSlackFile(&ssr.File, sfs.Endpoint.Id, sfs.Endpoint.UserId, sfs.Endpoint.Index)
	if err := file.IndexAsync(f, globals.FilesTopic, sfs.Endpoint.Index, true); err != nil {
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