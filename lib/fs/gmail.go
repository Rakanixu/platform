package fs

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/image"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

// GmailFs Gmail fyle system
type GmailFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

// NewGmailFsFromEndpoint constructor
func NewGmailFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GmailFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover attached files in google mail
func (gfs *GmailFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := gfs.getMessages(); err != nil {
			log.Println(err)
		}

		gfs.WalkRunning <- false
	}()

	return gfs.FilesChan, gfs.WalkRunning
}

// WalkUsers
func (gfs *GmailFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		gfs.WalkUsersRunning <- false
	}()

	return gfs.UsersChan, gfs.WalkUsersRunning
}

// WalkChannels
func (gfs *GmailFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		gfs.WalkChannelsRunning <- false
	}()

	return gfs.ChannelsChan, gfs.WalkChannelsRunning
}

// Token returns gmail user token
func (gfs *GmailFs) Token(c client.Client) string {
	return gfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (gfs *GmailFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (gfs *GmailFs) GetThumbnail(id string, c client.Client) (string, error) {
	return "", nil
}

// Create file in gmail (not implemented)
func (gfs *GmailFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	return gfs.FilesChan
}

// Delete (not implemented)
func (gfs *GmailFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	return gfs.FilesChan
}

// Update file
func (gfs *GmailFs) Update(req file_proto.ShareRequest) chan FileMsg {
	return gfs.FilesChan
}

// DownloadFile retrieves a file
func (gfs *GmailFs) DownloadFile(id string, opts ...string) (io.ReadCloser, error) {
	cfg := globals.NewGmailOauthConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	s, err := gmail.New(c)
	if err != nil {
		return nil, err
	}

	srv := gmail.NewUsersMessagesService(s)

	if len(opts) == 0 {
		return nil, errors.New("ERROR opts reuired (attachmentID)")
	}
	mpb, err := srv.Attachments.Get("me", id, opts[0]).Fields("data").Do()
	if err != nil {
		return nil, err
	}

	b, err := base64.URLEncoding.DecodeString(mpb.Data)
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(bytes.NewReader(b)), nil

}

// UploadFile uploads a file into google cloud storage
func (gfs *GmailFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, gfs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (gfs *GmailFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(gfs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (gfs *GmailFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(gfs.Endpoint.Index, "")
}

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
				if err := gfs.generateThumbnail(gf, v, vl, f.ID); err != nil {
					log.Println(err)
				}

				gfs.FilesChan <- NewFileMsg(f, nil)
			}
		}
	}

	return nil
}

// generateThumbnail downloads original picture, resize and uploads to Google storage
func (gfs *GmailFs) generateThumbnail(gf *gmailhelper.GmailFile, msg *gmail.Message, msgp *gmail.MessagePart, id string) error {
	if msgp.MimeType == globals.MIME_PNG || msgp.MimeType == globals.MIME_JPG || msgp.MimeType == globals.MIME_JPEG {
		pr, err := gfs.DownloadFile(msg.Id, msgp.Body.AttachmentId)
		if err != nil {
			return errors.New("ERROR downloading gmail file")
		}

		b, err := image.Thumbnail(pr, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for gmail file")
		}

		if err := gfs.UploadFile(b, id); err != nil {
			return errors.New("ERROR uploading thumbnail for gmail file: %s")
		}
	}

	return nil
}
