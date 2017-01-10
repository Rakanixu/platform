package fs

import (
	"errors"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/image"
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
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan file.File
	FileMetaChan        chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

//NewGoogleDriveFsFromEndpoint constructor
func NewGoogleDriveFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	return &GoogleDriveFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan file.File),
		FileMetaChan:        make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover files in google drive datasource
func (gfs *GoogleDriveFs) Walk() (chan file.File, chan bool, error) {
	go func() {
		if err := gfs.getFiles(); err != nil {
			log.Println(err)
		}

		gfs.WalkRunning <- false
	}()

	return gfs.FilesChan, gfs.WalkRunning, nil
}

// WalkUsers
func (gfs *GoogleDriveFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		gfs.WalkUsersRunning <- false
	}()

	return gfs.UsersChan, gfs.WalkUsersRunning
}

// WalkChannels
func (gfs *GoogleDriveFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		gfs.WalkChannelsRunning <- false
	}()

	return gfs.ChannelsChan, gfs.WalkChannelsRunning
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
func (gfs *GoogleDriveFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		f, err := srv.Files.Create(&drive.File{
			Name:     rq.FileName,
			MimeType: globals.GetMimeType(globals.GoogleDrive, rq.MimeType),
		}).Fields("*").Do()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		kfg := file.NewKazoupFileFromGoogleDriveFile(*f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if kfg == nil {
			gfs.FileMetaChan <- NewFileMsg(nil, errors.New("ERROR CreateFile gdrive is nil"))
			return
		}

		gfs.FileMetaChan <- NewFileMsg(kfg, nil)
	}()

	return gfs.FileMetaChan
}

// DeleteFile moves a google drive file to trash
func (gfs *GoogleDriveFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		// Trash file
		_, err = srv.Files.Update(rq.OriginalId, &drive.File{
			Trashed: true,
		}).Do()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		gfs.FileMetaChan <- NewFileMsg(
			&file.KazoupGoogleFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return gfs.FileMetaChan
}

// Update file
func (gfs *GoogleDriveFs) Update(req file_proto.ShareRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		if _, err := srv.Permissions.Create(req.OriginalId, &drive.Permission{
			Role:         "writer",
			Type:         "user",
			EmailAddress: req.DestinationId,
		}).Do(); err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		gf, err := srv.Files.Get(req.OriginalId).Fields("*").Do()
		if err != nil {
			gfs.FileMetaChan <- NewFileMsg(nil, err)
			return
		}

		kfg := file.NewKazoupFileFromGoogleDriveFile(*gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if kfg == nil {
			gfs.FileMetaChan <- NewFileMsg(nil, errors.New("ERROR ShareFile gdrive is nil"))
			return
		}

		gfs.FileMetaChan <- NewFileMsg(kfg, nil)
	}()

	return gfs.FileMetaChan
}

// DownloadFile retrieves a file from google drive
func (gfs *GoogleDriveFs) DownloadFile(id string, opts ...string) (io.ReadCloser, error) {
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
		f := file.NewKazoupFileFromGoogleDriveFile(*v, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if f != nil {
			if err := gfs.generateThumbnail(v, f.ID); err != nil {
				log.Println(err)
			}

			gfs.FilesChan <- f
		}
	}

	return nil
}

func (gfs *GoogleDriveFs) generateThumbnail(f *drive.File, id string) error {
	c := categories.GetDocType("." + f.FullFileExtension)
	if len(f.FullFileExtension) == 0 {
		c = categories.GetDocType(f.MimeType)
	}

	if c == globals.CATEGORY_PICTURE {
		rc, err := gfs.DownloadFile(f.Id)
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
