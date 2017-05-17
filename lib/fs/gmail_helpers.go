package fs

import (
	"fmt"
	"github.com/cenkalti/backoff"
	cs "github.com/kazoup/platform/lib/cloudstorage"
	"github.com/kazoup/platform/lib/cloudvision"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/image"
	"github.com/kazoup/platform/lib/objectstorage"
	sttlib "github.com/kazoup/platform/lib/speechtotext"
	"github.com/kazoup/platform/lib/tika"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// getMessages discover files (attachments)
func (gfs *GmailFs) getMessages() error {
	cfg := globals.NewGmailOauthConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	s, err := gmail.New(c)
	if err != nil {
		return err
	}

	srv := gmail.NewUsersMessagesService(s)
	srvCall := srv.List("me") // Token authenticate user
	msgBdy, err := srvCall.Q("has:attachment").Fields("messages,nextPageToken,resultSizeEstimate").Do()
	if err != nil {
		return err
	}

	if len(msgBdy.Messages) > 0 {
		if err := gfs.pushMessagesToChanForPage(s, msgBdy.Messages); err != nil {
			return err
		}
	}

	if len(msgBdy.NextPageToken) > 0 {
		if err := gfs.getNextPage(s, msgBdy.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// getNextPage allows pagination for discovering files
func (gfs *GmailFs) getNextPage(s *gmail.Service, nextPageToken string) error {
	srv := gmail.NewUsersMessagesService(s)
	r, err := srv.List("me").PageToken(nextPageToken).Fields("messages,nextPageToken,resultSizeEstimate").Do()
	if err != nil {
		return err
	}

	if len(r.Messages) > 0 {
		if err := gfs.pushMessagesToChanForPage(s, r.Messages); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(s, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// pushMessagesToChanForPage push discovered files to broker
func (gfs *GmailFs) pushMessagesToChanForPage(s *gmail.Service, msgs []*gmail.Message) error {
	srv := gmail.NewUsersMessagesService(s)

	for _, v := range msgs {
		// Available fields
		// historyId,id,internalDate,labelIds,payload,raw,sizeEstimate,snippet,threadId
		msgBdy, err := srv.Get("me", v.Id).Fields("id,internalDate,payload,sizeEstimate").Do()
		if err != nil {
			return err
		}

		// Iterate over all attachments
		for _, vl := range msgBdy.Payload.Parts {
			gf := &gmailhelper.GmailFile{
				Id:           fmt.Sprintf("%s%s", msgBdy.Id, vl.PartId),
				MessageId:    msgBdy.Id,
				AttachmentId: vl.Body.AttachmentId,
				Extension:    "None",
				InternalDate: msgBdy.InternalDate,
				SizeEstimate: msgBdy.SizeEstimate,
				Name:         vl.Filename,
				MimeType:     vl.MimeType,
			}

			ext := strings.Split(strings.Replace(vl.Filename, " ", "-", 1), ".")
			gf.Extension = ext[len(ext)-1]

			f := file.NewKazoupFileFromGmailFile(*gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Url, gfs.Endpoint.Index)
			// Constructor will return nil when the attachment has no name
			// When an attachment has no name, attachment use to be a marketing image
			if f != nil {
				gfs.FilesChan <- NewFileMsg(f, nil)
			}
		}
	}

	return nil
}

// processImage, cloud vision
func (gfs *GmailFs) processImage(f *file.KazoupGmailFile) (file.File, error) {
	// Downloads from gmail, see connector
	gmcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.Gmail)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = gmcs.Download(f.MessageId, f.AttachmentId)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	// Resize to optimal size for cloud vision API
	cvrd, err := image.Thumbnail(rc, globals.CLOUD_VISION_IMG_WIDTH)
	if err != nil {
		return nil, err
	}

	if f.Tags, err = cloudvision.Tag(ioutil.NopCloser(cvrd)); err != nil {
		return nil, err
	}

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			TagsTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.TagsTimestamp = &n
	}

	return f, nil
}

// enrichFile sends the original file to tika and enrich KazoupGmailFile with Tika interface
func (gfs *GmailFs) processDocument(f *file.KazoupGmailFile) (file.File, error) {
	// Download file from Gmail, so connector is globals.Gmail
	gcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.Gmail)
	if err != nil {
		return nil, err
	}

	rc, err := gcs.Download(f.MessageId, f.AttachmentId)
	if err != nil {
		return nil, err
	}

	t, err := tika.ExtractPlainContent(rc)
	if err != nil {
		return nil, err
	}

	f.Content = t.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ContentTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ContentTimestamp = &n
	}

	return f, nil
}

// processAudio uploads audio file to GCS and runs async speech to text over it
func (gfs *GmailFs) processAudio(f *file.KazoupGmailFile) (file.File, error) {
	// Download file from Box, so connector is globals.Box
	gmcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.Gmail)
	if err != nil {
		return nil, err
	}

	rc, err := gmcs.Download(f.MessageId, f.AttachmentId)
	if err != nil {
		return nil, err
	}

	if err := objectstorage.Upload(rc, globals.AUDIO_BUCKET, f.ID); err != nil {
		return nil, err
	}

	stt, err := sttlib.AsyncContent(fmt.Sprintf("gs://%s/%s", globals.AUDIO_BUCKET, f.ID))
	if err != nil {
		return nil, err
	}

	f.Content = stt.Content()

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			AudioTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.AudioTimestamp = &n
	}

	return f, nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (gfs *GmailFs) processThumbnail(f *file.KazoupGmailFile) (file.File, error) {
	// Downloads from gmail, see connector
	gmcs, err := cs.NewCloudStorageFromEndpoint(gfs.Endpoint, globals.Gmail)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser

	if err := backoff.Retry(func() error {
		rc, err = gmcs.Download(f.MessageId, f.AttachmentId)
		if err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff()); err != nil {
		return nil, err
	}
	defer rc.Close()

	backoff.Retry(func() error {
		b, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			// Skip retry
			return nil
		}

		if err := objectstorage.Upload(ioutil.NopCloser(b), gfs.Endpoint.Index, f.ID); err != nil {
			return err
		}

		return nil
	}, backoff.NewExponentialBackOff())

	n := time.Now()
	if f.OptsKazoupFile == nil {
		f.OptsKazoupFile = &file.OptsKazoupFile{
			ThumbnailTimestamp: &n,
		}
	} else {
		f.OptsKazoupFile.ThumbnailTimestamp = &n
	}

	return f, nil
}
