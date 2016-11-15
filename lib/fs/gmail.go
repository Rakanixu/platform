package fs

import (
	"encoding/base64"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
	"log"
	"strings"
	"time"
)

// GmailFs Gmail fyle system
type GmailFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

// NewGmailFsFromEndpoint constructor
func NewGmailFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GmailFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

// List returns 2 channels, for files and state. Discover attached files in google mail
func (gfs *GmailFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getMessages(); err != nil {
			log.Println(err)
		}

		gfs.Running <- false
	}()

	return gfs.FilesChan, gfs.Running, nil
}

// Token returns gmail user token
func (gfs *GmailFs) Token() string {
	return gfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (gfs *GmailFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (gfs *GmailFs) GetThumbnail(id string) (string, error) {
	return "", nil
}

// CreateFile belongs to Fs interface
func (gfs *GmailFs) CreateFile(rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	return &file_proto.CreateResponse{}, nil
}

// DeleteFile deletes a email and therefore its attachments. Be careful.
func (gfs *GmailFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	return &file_proto.DeleteResponse{}, nil
}

// ShareFile
func (gfs *GmailFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	return "", nil
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
			}

			if vl.MimeType == globals.MIME_PNG || vl.MimeType == globals.MIME_JPG || vl.MimeType == globals.MIME_JPEG {
				mpb, err := srv.Attachments.Get("me", v.Id, vl.Body.AttachmentId).Fields("data").Do()
				if err != nil {
					return err
				}

				b, err := base64.URLEncoding.DecodeString(mpb.Data)
				if err != nil {
					return err
				}

				gf.Base64 = string(base64.StdEncoding.EncodeToString(b))
			}

			ext := strings.Split(strings.Replace(vl.Filename, " ", "-", 1), ".")
			gf.Extension = ext[len(ext)-1]

			f := file.NewKazoupFileFromGmailFile(gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Url, gfs.Endpoint.Index)
			// Constructor will return nil when the attachment has no name
			// When an attachment has no name, attachment use to be a marketing image
			if f != nil {
				gfs.FilesChan <- f
			}
		}
	}

	return nil
}
