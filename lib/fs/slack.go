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
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// SlackFs slack file system
type SlackFs struct {
	Endpoint     *datasource_proto.Endpoint
	Running      chan bool
	FilesChan    chan file.File
	FileMetaChan chan FileMeta
}

// NewSlackFsFromEndpoint constructor
func NewSlackFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &SlackFs{
		Endpoint:     e,
		Running:      make(chan bool, 1),
		FilesChan:    make(chan file.File),
		FileMetaChan: make(chan FileMeta),
	}
}

// List returns 2 channels, for files and state. Discover files in slack datasource
func (sfs *SlackFs) List(c client.Client) (chan file.File, chan bool, error) {
	go func() {
		if err := sfs.getUsers(c); err != nil {
			log.Println(err)
		}

		if err := sfs.getChannels(c); err != nil {
			log.Println(err)
		}

		if err := sfs.getFiles(c, 1); err != nil {
			log.Println(err)
		}
		// Slack scan finished
		sfs.Running <- false
	}()

	return sfs.FilesChan, sfs.Running, nil
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
func (sfs *SlackFs) Create(rq file_proto.CreateRequest) chan FileMeta {
	return sfs.FileMetaChan
}

// DeleteFile deletes a slack file
func (sfs *SlackFs) Delete(rq file_proto.DeleteRequest) chan FileMeta {
	return sfs.FileMetaChan
}

// ShareFile sets a PermalinkPublic available, so everyone with URL has access to the slack file
func (sfs *SlackFs) Update(req file_proto.ShareRequest) chan FileMeta {
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
	return sfs.FileMetaChan
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
func (sfs *SlackFs) getUsers(cl client.Client) error {
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

		if err := cl.Publish(context.Background(), cl.NewPublication(globals.SlackUsersTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

// getChannels retrieves channels from slack team
func (sfs *SlackFs) getChannels(cl client.Client) error {
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

		if err := cl.Publish(context.Background(), cl.NewPublication(globals.SlackChannelsTopic, msg)); err != nil {
			return err
		}
	}

	return nil
}

// getFiles discover slack files for user and push them to broker
func (sfs *SlackFs) getFiles(cl client.Client, page int) error {
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

		if err := sfs.generateThumbnail(cl, v, f.ID); err != nil {
			log.Println(err)
		}

		sfs.FilesChan <- f
	}

	if filesRsp.Paging.Pages >= page {
		sfs.getFiles(cl, page+1)
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (sfs *SlackFs) generateThumbnail(c client.Client, sf slack.SlackFile, id string) error {
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
