package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
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
			MimeType: utils.GetMimeType(globals.GoogleDrive, rq.MimeType),
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
