package fs

import (
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	file "github.com/kazoup/platform/lib/file"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"log"
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

// Enrich
func (gfs *GmailFs) Enrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		//bfs.processImage()
		//bfs.processDocument()
	}()

	return gfs.FilesChan
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
