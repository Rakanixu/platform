package fs

import (
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"log"
	"time"
)

type GoogleDriveFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

func NewGoogleDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GoogleDriveFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

func (gfs *GoogleDriveFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getFiles(); err != nil {
			log.Println(err)
		}

		gfs.Running <- false
	}()

	return gfs.FilesChan, gfs.Running, nil
}

func (gfs *GoogleDriveFs) Token() string {
	return gfs.Endpoint.Token.AccessToken
}

func (gfs *GoogleDriveFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

func (gfs *GoogleDriveFs) GetThumbnail(id string) (string, error) {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	srv, err := drive.New(c)
	if err != nil {
		return "", err
	}
	r, err := srv.Files.Get(id).Fields("thumbnailLink").Do()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%ss700", r.ThumbnailLink[:len(r.ThumbnailLink)-4]), nil
}

func (gfs *GoogleDriveFs) CreateFile(fileType string) (string, error) {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	srv, err := drive.New(c)
	if err != nil {
		return "", err
	}

	f, err := srv.Files.Create(&drive.File{
		MimeType: globals.GetMimeType(globals.GoogleDrive, fileType),
	}).Fields("*").Do()
	if err != nil {
		return "", err
	}

	kfg := file.NewKazoupFileFromGoogleDriveFile(f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
	if err := file.IndexAsync(kfg, globals.FilesTopic, gfs.Endpoint.Index); err != nil {
		return "", err
	}

	return kfg.GetURL(), nil
}

func (gfs *GoogleDriveFs) getFiles() error {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	srv, err := drive.New(c)
	if err != nil {
		return err
	}

	r, err := srv.Files.List().PageSize(100).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

func (gfs *GoogleDriveFs) getNextPage(srv *drive.Service, nextPageToken string) error {
	r, err := srv.Files.List().PageToken(nextPageToken).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

func (gfs *GoogleDriveFs) pushFilesToChanForPage(files []*drive.File) error {
	for _, v := range files {
		f := file.NewKazoupFileFromGoogleDriveFile(v, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)

		gfs.FilesChan <- f
	}

	return nil
}
