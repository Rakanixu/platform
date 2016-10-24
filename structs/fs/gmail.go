package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	gmail "google.golang.org/api/gmail/v1"
	"log"
	"time"
)

type GmailFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

func NewGmailFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GmailFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

func (gfs *GmailFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getMessages(); err != nil {
			log.Println(err)
		}

		gfs.Running <- false
	}()

	return gfs.FilesChan, gfs.Running, nil
}

func (gfs *GmailFs) Token() string {
	return gfs.Endpoint.Token.AccessToken
}

func (gfs *GmailFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

func (gfs *GmailFs) GetThumbnail(id string) (string, error) {
	return "", nil
}

func (gfs *GmailFs) CreateFile(fileType string) (string, error) {
	return "", nil
}

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

func (gfs *GmailFs) pushMessagesToChanForPage(s *gmail.Service, msgs []*gmail.Message) error {
	srv := gmail.NewUsersMessagesService(s)

	for _, v := range msgs {
		srvCall := srv.Get("me", v.Id)
		msgBdy, err := srvCall.Fields("historyId,id,internalDate,labelIds,payload,raw,sizeEstimate,snippet,threadId").Do()
		if err != nil {
			return err
		}

		f := file.NewKazoupFileFromGmailFile(msgBdy, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Url, gfs.Endpoint.Index)
		// Constructor will return nil when the attachment has no name
		// When an attachment has no name, attachment use to be a marketing image
		if f != nil {
			gfs.FilesChan <- f
		}
	}

	return nil
}
