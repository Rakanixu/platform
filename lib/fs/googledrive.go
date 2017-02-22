package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"google.golang.org/api/drive/v3"
	"log"
)

// GoogleDriveFs is the google drive file system struct
type GoogleDriveFs struct {
	Endpoint            *datasource_proto.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
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
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover files in google drive datasource
func (gfs *GoogleDriveFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := gfs.getFiles(); err != nil {
			log.Println(err)
		}

		gfs.WalkRunning <- false
	}()

	return gfs.FilesChan, gfs.WalkRunning
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

// Enrich
func (gfs *GoogleDriveFs) Enrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		var err error

		_, ok := f.(*file.KazoupGoogleFile)
		if !ok {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("Error enriching file"))
			return
		}

		// OptsKazoupFile.ContentTimestamp and
		// OptsKazoupFile.TagsTimestamp are not defined,
		// Content was never extracted before
		process := struct {
			Picture  bool
			Document bool
		}{
			Picture:  false,
			Document: false,
		}
		if f.(*file.KazoupGoogleFile).OptsKazoupFile == nil {
			process.Picture = true
			process.Document = true
		} else {
			process.Picture = f.(*file.KazoupGoogleFile).OptsKazoupFile.TagsTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
			process.Document = f.(*file.KazoupGoogleFile).OptsKazoupFile.ContentTimestamp.Before(f.(*file.KazoupGoogleFile).Modified)
		}

		/*		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_PICTURE && process.Picture {
				f, err = gfs.processImage(gcs, f.(*file.KazoupGoogleFile))
				if err != nil {
					gfs.FilesChan <- NewFileMsg(nil, err)
					return
				}
			}*/

		if f.(*file.KazoupGoogleFile).Category == globals.CATEGORY_DOCUMENT && process.Document {
			f, err = gfs.processDocument(f.(*file.KazoupGoogleFile))
			if err != nil {
				gfs.FilesChan <- NewFileMsg(nil, err)
				return
			}
		}

		gfs.FilesChan <- NewFileMsg(f, err)
	}()

	return gfs.FilesChan
}

// CreateFile creates a google file and index it on Elastic Search
func (gfs *GoogleDriveFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		f, err := srv.Files.Create(&drive.File{
			Name:     rq.FileName,
			MimeType: globals.GetMimeType(globals.GoogleDrive, rq.MimeType),
		}).Fields("*").Do()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		kfg := file.NewKazoupFileFromGoogleDriveFile(*f, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if kfg == nil {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("ERROR CreateFile gdrive is nil"))
			return
		}

		gfs.FilesChan <- NewFileMsg(kfg, nil)
	}()

	return gfs.FilesChan
}

// DeleteFile moves a google drive file to trash
func (gfs *GoogleDriveFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		// Trash file
		_, err = srv.Files.Update(rq.OriginalId, &drive.File{
			Trashed: true,
		}).Do()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		// Return deleted file. This file only stores the id
		// Avoid read from DB
		gfs.FilesChan <- NewFileMsg(
			&file.KazoupGoogleFile{
				file.KazoupFile{
					ID: rq.FileId,
				},
				nil,
			},
			nil,
		)
	}()

	return gfs.FilesChan
}

// Update file
func (gfs *GoogleDriveFs) Update(req file_proto.ShareRequest) chan FileMsg {
	go func() {
		srv, err := gfs.getDriveService()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		if _, err := srv.Permissions.Create(req.OriginalId, &drive.Permission{
			Role:         "writer",
			Type:         "user",
			EmailAddress: req.DestinationId,
		}).Do(); err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		gf, err := srv.Files.Get(req.OriginalId).Fields("*").Do()
		if err != nil {
			gfs.FilesChan <- NewFileMsg(nil, err)
			return
		}

		kfg := file.NewKazoupFileFromGoogleDriveFile(*gf, gfs.Endpoint.Id, gfs.Endpoint.UserId, gfs.Endpoint.Index)
		if kfg == nil {
			gfs.FilesChan <- NewFileMsg(nil, errors.New("ERROR ShareFile gdrive is nil"))
			return
		}

		gfs.FilesChan <- NewFileMsg(kfg, nil)
	}()

	return gfs.FilesChan
}
