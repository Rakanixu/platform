package fs

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
)

type MockFs struct {
	Endpoint            *proto_datasource.Endpoint
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

func NewMockFsFromEndpoint(e *proto_datasource.Endpoint) Fs {
	return &MockFs{
		Endpoint:            e,
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

func (mfs *MockFs) Walk() (chan FileMsg, chan bool) {
	// Set not running after start
	go func() {
		mfs.WalkRunning <- false
	}()

	return mfs.FilesChan, mfs.WalkRunning
}

func (mfs *MockFs) WalkUsers() (chan UserMsg, chan bool) {
	// Set not running after start
	go func() {
		mfs.WalkUsersRunning <- false
	}()

	return mfs.UsersChan, mfs.WalkUsersRunning
}

func (mfs *MockFs) WalkChannels() (chan ChannelMsg, chan bool) {
	// Set not running after start
	go func() {
		mfs.WalkChannelsRunning <- false
	}()

	return mfs.ChannelsChan, mfs.WalkChannelsRunning
}

func (mfs *MockFs) Create(rq proto_file.CreateRequest) chan FileMsg {
	go func() {
		// Create Mock file
		kfo := file.NewKazoupFileFromMockFile()

		mfs.FilesChan <- NewFileMsg(kfo, nil)
	}()

	return mfs.FilesChan
}

func (mfs *MockFs) Delete(rq proto_file.DeleteRequest) chan FileMsg {
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}

func (mfs *MockFs) Update(req proto_file.ShareRequest) chan FileMsg {
	go func() {
		mfs.FilesChan <- NewFileMsg(file.NewKazoupFileFromMockFile(), nil)
	}()

	return mfs.FilesChan
}
