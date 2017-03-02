package fs

import (
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/slack"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
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

		sfs.FilesChan <- NewFileMsg(f, nil)
	}

	if filesRsp.Paging.Pages >= page {
		sfs.getFiles(page + 1)
	}

	return nil
}

// processImage, cloud vision processing
func (sfs *SlackFs) processImage(f *file.KazoupSlackFile) (file.File, error) {
	// Download file from Slack, so connector is globals.Slack
	scs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.Slack)

	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = scs.Download(f.Original.URLPrivateDownload)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		log.Println("ERROR DOWNLOADING FILE", err)
		return nil, err
	}
	defer rc.Close()

	// Resize to optimal size for cloud vision API
	cvrd, err := image.Thumbnail(rc, globals.CLOUD_VISION_IMG_WIDTH)
	if err != nil {
		log.Println("CLOUD VISION ERROR", err)
		return nil, err
	}

	if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
		log.Println("CLOUD VISION ERROR", err)
		return nil, err
	}

	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			TagsTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.TagsTimestamp = time.Now()
	}

	return f, nil
}

// enrichFile sends the original file to tika and enrich KazoupOneDriveFile with Tika interface
func (sfs *SlackFs) processDocument(f *file.KazoupSlackFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.OneDrive
	scs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.Slack)
	if err != nil {
		return nil, err
	}

	rc, err := scs.Download(f.Original.URLPrivateDownload)
	if err != nil {
		return nil, err
	}

	t, err := tika.ExtractContent(rc)
	if err != nil {
		return nil, err
	}

	f.Content = t.Content()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ContentTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.ContentTimestamp = time.Now()
	}

	return f, nil
}

// processAudio uploads audio file to GCS and runs async speech to text over it
func (sfs *SlackFs) processAudio(gcs *gcslib.GoogleCloudStorage, f *file.KazoupSlackFile) (file.File, error) {
	// Download file from GoogleDrive, so connector is globals.OneDrive
	scs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.Slack)
	if err != nil {
		return nil, err
	}

	rc, err := scs.Download(f.Original.URLPrivateDownload)
	if err != nil {
		return nil, err
	}

	if err := gcs.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

	stt, err := sttlib.AsyncContent(fmt.Sprintf("gs://%s/%s", globals.AUDIO_BUCKET, f.ID))
	if err != nil {
		return nil, err
	}

	f.Content = stt.Content()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			AudioTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.AudioTimestamp = time.Now()
	}

	return f, nil
}

// processThumbnail downloads original picture, resize and uploads to Google storage
func (sfs *SlackFs) processThumbnail(gcs *gcslib.GoogleCloudStorage, f *file.KazoupSlackFile) (file.File, error) {
	// Download file from Slack, so connector is globals.Slack
	scs, err := cs.NewCloudStorageFromEndpoint(sfs.Endpoint, globals.Slack)

	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = scs.Download(f.Original.URLPrivateDownload)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		log.Println("ERROR DOWNLOADING FILE", err)
		return nil, err
	}
	defer rc.Close()

	backoff.Retry(func() error {
		b, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			log.Println("THUMNAIL GENERATION ERROR, SKIPPING", err)
			// Skip retry
			return nil
		}

		if err := gcs.Upload(ioutil.NopCloser(b), sfs.Endpoint.Index, f.ID); err != nil {
			log.Println("THUMNAIL UPLOAD ERROR", err)
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())

	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ThumbnailTimestamp: time.Now(),
		}
	} else {
		f.OptsKazoupFile.ThumbnailTimestamp = time.Now()
	}

	return f, nil
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
