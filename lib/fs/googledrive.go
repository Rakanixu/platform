package fs

import (
	"errors"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"time"
)

// GoogleDriveFs is the google drive file system struct
type GoogleDriveFs struct {
	Endpoint     *datasource_proto.Endpoint
	Running      chan bool
	FilesChan    chan file.File
	FileMetaChan chan FileMeta
}

//NewGoogleDriveFsFromEndpoint constructor
func NewGoogleDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GoogleDriveFs{
		Endpoint:     e,
		Running:      make(chan bool, 1),
		FilesChan:    make(chan file.File),
		FileMetaChan: make(chan FileMeta),
	}
}

// List returns 2 channels, for files and state. Discover files in google drive datasource
func (gfs *GoogleDriveFs) List(c client.Client) (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getFiles(c); err != nil {
			log.Println(err)
		}

		gfs.Running <- false
	}()

	return gfs.FilesChan, gfs.Running, nil
}

// Token returns google drive user token
func (gfs *GoogleDriveFs) Token(c client.Client) string {
	return gfs.Endpoint.Token.AccessToken
}

// GetDatasourceId returns datasource ID
func (gfs *GoogleDriveFs) GetDatasourceId() string {
	return gfs.Endpoint.Id
}

// GetThumbnail returns a URI pointing to a thumbnail
func (gfs *GoogleDriveFs) GetThumbnail(id string, c client.Client) (string, error) {
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
func (gfs *GoogleDriveFs) Create(rq file_proto.CreateRequest) chan FileMeta {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		f, err := srv.Files.Create(&drive.File{
			Name:     rq.FileName,
			MimeType: globals.GetMimeType(globals.GoogleDrive, rq.MimeType),
		}).Fields("*").Do()
		if err != nil {
			gfs.FileMetaChan <- NewFileMeta(nil, err)
			return
		}

		kfg := file.NewKazoupFileFromGoogleDriveFile(*f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if kfg == nil {
			gfs.FileMetaChan <- NewFileMeta(nil, errors.New("ERROR CreateFile gdrive is nil"))
			return
		}

		gfs.FileMetaChan <- NewFileMeta(kfg, nil)
	}()

	return gfs.FileMetaChan
}

// DeleteFile moves a google drive file to trash
func (gfs *GoogleDriveFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return nil, err
	}

	// Trash file
	_, err = srv.Files.Update(rq.OriginalId, &drive.File{
		Trashed: true,
	}).Do()
	if err != nil {
		return nil, err
	}

	// Delete file from index
	req := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Delete",
		&db_proto.DeleteRequest{
			Index: gfs.Endpoint.Index,
			Type:  globals.FileType,
			Id:    rq.FileId,
		},
	)
	rsp := &db_proto.DeleteResponse{}

	if err := c.Call(ctx, req, rsp); err != nil {
		return nil, err
	}

	// Publish notification topic, let client know when to refresh itself
	if err := c.Publish(globals.NewSystemContext(), c.NewPublication(globals.NotificationTopic, &notification_proto.NotificationMessage{
		Method: globals.NOTIFY_REFRESH_SEARCH,
		UserId: gfs.Endpoint.UserId,
	})); err != nil {
		log.Print("Publishing (notify file) error %s", err)
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

	kfg := file.NewKazoupFileFromGoogleDriveFile(*gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
	if kfg == nil {
		return "", errors.New("ERROR ShareFile gdrive is nil")
	}

	if err := file.IndexAsync(c, kfg, globals.FilesTopic, gfs.Endpoint.Index, true); err != nil {
		return "", err
	}

	return "", nil
}

// DownloadFile retrieves a file from google drive
func (gfs *GoogleDriveFs) DownloadFile(id string, c client.Client, opts ...string) (io.ReadCloser, error) {
	srv, err := gfs.getDriveService()
	if err != nil {
		return nil, err
	}

	res, err := srv.Files.Get(id).Download()
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// UploadFile uploads a file into google cloud storage
func (gfs *GoogleDriveFs) UploadFile(file io.Reader, fId string) error {
	return UploadFile(file, gfs.Endpoint.Index, fId)
}

// SignedObjectStorageURL returns a temporary link to a resource in GC storage
func (gfs *GoogleDriveFs) SignedObjectStorageURL(objName string) (string, error) {
	return SignedObjectStorageURL(gfs.Endpoint.Index, objName)
}

// DeleteFilesFromIndex removes files from GC storage
func (gfs *GoogleDriveFs) DeleteIndexBucketFromGCS() error {
	return DeleteBucket(gfs.Endpoint.Index, "")
}

// getFiles discover all files in google drive account
func (gfs *GoogleDriveFs) getFiles(c client.Client) error {
	srv, err := gfs.getDriveService()
	if err != nil {
		return err
	}

	r, err := srv.Files.List().PageSize(100).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(c, r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(c, srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// getNextPage allows pagination while discovering files
func (gfs *GoogleDriveFs) getNextPage(c client.Client, srv *drive.Service, nextPageToken string) error {
	r, err := srv.Files.List().PageToken(nextPageToken).Fields("files,kind,nextPageToken").Do()
	if err != nil {
		return err
	}

	if len(r.Files) > 0 {
		if err := gfs.pushFilesToChanForPage(c, r.Files); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gfs.getNextPage(c, srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}

// pushFilesToChanForPage sends discovered files to the file system channel
func (gfs *GoogleDriveFs) pushFilesToChanForPage(c client.Client, files []*drive.File) error {
	for _, v := range files {
		f := file.NewKazoupFileFromGoogleDriveFile(*v, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if f != nil {
			if err := gfs.generateThumbnail(c, v, f.ID); err != nil {
				log.Println(err)
			}

			gfs.FilesChan <- f
		}
	}

	return nil
}

func (gfs *GoogleDriveFs) generateThumbnail(cl client.Client, f *drive.File, id string) error {
	c := categories.GetDocType("." + f.FullFileExtension)
	if len(f.FullFileExtension) == 0 {
		c = categories.GetDocType(f.MimeType)
	}

	if c == globals.CATEGORY_PICTURE {
		rc, err := gfs.DownloadFile(f.Id, cl)
		if err != nil {
			return errors.New("ERROR downloading googledrive file")
		}

		rd, err := image.Thumbnail(rc, globals.THUMBNAIL_WIDTH)
		if err != nil {
			return errors.New("ERROR generating thumbnail for googledrive file")
		}

		if err := gfs.UploadFile(rd, id); err != nil {
			return errors.New("ERROR uploading thumbnail for googledrive file")
		}
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
