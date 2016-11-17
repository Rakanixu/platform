package fs

import (
	"encoding/json"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"time"
)

// GoogleDriveFs is the google drive file system struct
type GoogleDriveFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	FilesChan chan file.File
}

//NewGoogleDriveFsFromEndpoint constructor
func NewGoogleDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GoogleDriveFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		FilesChan: make(chan file.File),
	}
}

// List returns 2 channels, for files and state. Discover files in google drive datasource
func (gfs *GoogleDriveFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getFiles(); err != nil {
			log.Println(err)
		}

		gfs.Running <- false
	}()

	return gfs.FilesChan, gfs.Running, nil
}

// Token returns google drive user token
func (gfs *GoogleDriveFs) Token() string {
	return gfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (gfs *GoogleDriveFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (gfs *GoogleDriveFs) GetThumbnail(id string) (string, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return "", err
	}
	r, err := srv.Files.Get(id).Fields("thumbnailLink").Do()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%ss700", r.ThumbnailLink[:len(r.ThumbnailLink)-4]), nil
}

// CreateFile creates a google file and index it on Elastic Search
func (gfs *GoogleDriveFs) CreateFile(rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return nil, err
	}

	f, err := srv.Files.Create(&drive.File{
		Name:     rq.FileName,
		MimeType: globals.GetMimeType(globals.GoogleDrive, rq.MimeType),
	}).Fields("*").Do()
	if err != nil {
		return nil, err
	}

	kfg := file.NewKazoupFileFromGoogleDriveFile(f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
	if err := file.IndexAsync(kfg, globals.FilesTopic, gfs.Endpoint.Index, true); err != nil {
		return nil, err
	}

	b, err := json.Marshal(kfg)
	if err != nil {
		return nil, err
	}

	return &file_proto.CreateResponse{
		DocUrl: kfg.GetURL(),
		Data:   string(b),
	}, nil
}

// DeleteFile moves a google drive file to trash
func (gfs *GoogleDriveFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return nil, err
	}

	// Trash file
	f, err := srv.Files.Update(rq.OriginalId, &drive.File{
		Trashed: true,
	}).Fields("*").Do()
	if err != nil {
		return nil, err
	}

	// Reindex file
	kfg := file.NewKazoupFileFromGoogleDriveFile(f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
	if err := file.IndexAsync(kfg, globals.FilesTopic, gfs.Endpoint.Index, true); err != nil {
		return nil, err
	}

	return &file_proto.DeleteResponse{}, nil
}

// ShareFile
func (gfs *GoogleDriveFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return "", err
	}

	if _, err := srv.Permissions.Create(req.OriginalId, &drive.Permission{
		Role:         "writer",
		Type:         "user",
		EmailAddress: req.DestinationId,
	}).Do(); err != nil {
		return "", err
	}

	gf, err := srv.Files.Get(req.OriginalId).Fields("*").Do()
	if err != nil {
		return "", err
	}

	kfg := file.NewKazoupFileFromGoogleDriveFile(gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
	if err := file.IndexAsync(kfg, globals.FilesTopic, gfs.Endpoint.Index, true); err != nil {
		return "", err
	}

	return "", nil
}

// DownloadFile retrieves a file
func (gfs *GoogleDriveFs) DownloadFile(id string) ([]byte, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return nil, err
	}

	res, err := srv.Files.Get(id).Download()
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

// getFiles discover all files in google drive account
func (gfs *GoogleDriveFs) getFiles() error {
	srv, err := gfs.getDriveService()
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

// getNextPage allows pagination while discovering files
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

// pushFilesToChanForPage sends discovered files to the file system channel
func (gfs *GoogleDriveFs) pushFilesToChanForPage(files []*drive.File) error {
	for _, v := range files {
		c := categories.GetDocType("." + v.FullFileExtension)
		if len(v.FullFileExtension) == 0 {
			c = categories.GetDocType(v.MimeType)
		}

		b64 := ""
		if c == globals.CATEGORY_PICTURE {
			b, err := gfs.DownloadFile(v.Id)
			if err != nil {
				return err
			}

			b64, err = FileToBase64(b)
			if err != nil {
				return err
			}
		}

		f := file.NewKazoupFileFromGoogleDriveFile(v, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index, b64)

		gfs.FilesChan <- f
	}

	return nil
}

// getDriveService return a google drive service instance
func (gfs *GoogleDriveFs) getDriveService() (*drive.Service, error) {
	cfg := globals.NewGoogleOautConfig()
	c := cfg.Client(context.Background(), &oauth2.Token{
		AccessToken:  gfs.Endpoint.Token.AccessToken,
		TokenType:    gfs.Endpoint.Token.TokenType,
		RefreshToken: gfs.Endpoint.Token.RefreshToken,
		Expiry:       time.Unix(gfs.Endpoint.Token.Expiry, 0),
	})

	return drive.New(c)
}
