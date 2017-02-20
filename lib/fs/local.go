package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_config_proto "github.com/kazoup/platform/db/srv/proto/config"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/local"
	"golang.org/x/net/context"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LocalFs struct
type LocalFs struct {
	Endpoint            *datasource_proto.Endpoint
	RootPath            string
	WalkRunning         chan bool
	WalkUsersRunning    chan bool
	WalkChannelsRunning chan bool
	FilesChan           chan FileMsg
	UsersChan           chan UserMsg
	ChannelsChan        chan ChannelMsg
}

// NewLocalFsFromEndpoint constructor
func NewLocalFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	url := strings.Split(e.Url, "://")

	return &LocalFs{
		Endpoint:            e,
		RootPath:            url[1],
		WalkRunning:         make(chan bool, 1),
		WalkUsersRunning:    make(chan bool, 1),
		WalkChannelsRunning: make(chan bool, 1),
		FilesChan:           make(chan FileMsg),
		UsersChan:           make(chan UserMsg),
		ChannelsChan:        make(chan ChannelMsg),
	}
}

// Walk returns 2 channels, for files and state. Discover local files
func (lfs *LocalFs) Walk() (chan FileMsg, chan bool) {
	go func() {
		if err := lfs.walkDatasourceParents(); err != nil {
			log.Println("ERROR", err)
		}

		if err := filepath.Walk(lfs.RootPath, lfs.walkHandler()); err != nil {
			log.Println("ERROR", err)
		}
		lfs.WalkRunning <- false
	}()

	return lfs.FilesChan, lfs.WalkRunning
}

// WalUsers
func (lfs *LocalFs) WalkUsers() (chan UserMsg, chan bool) {
	go func() {
		// We can discover FS users, like root in Unix like
		lfs.WalkUsersRunning <- false
	}()

	return lfs.UsersChan, lfs.WalkUsersRunning
}

// WalkChannels
func (lfs *LocalFs) WalkChannels() (chan ChannelMsg, chan bool) {
	go func() {
		lfs.WalkChannelsRunning <- false
	}()

	return lfs.ChannelsChan, lfs.WalkChannelsRunning
}

// Enrich
func (lfs *LocalFs) Enrich(f file.File, gcs *gcslib.GoogleCloudStorage) chan FileMsg {
	go func() {
		//bfs.processImage()
		//bfs.processDocument()
	}()

	return lfs.FilesChan
}

// Create file (not implemented)
func (lfs *LocalFs) Create(rq file_proto.CreateRequest) chan FileMsg {
	return lfs.FilesChan
}

// DeleteFile deletes a local file
func (lfs *LocalFs) Delete(rq file_proto.DeleteRequest) chan FileMsg {
	return lfs.FilesChan
}

// Update file
func (lfs *LocalFs) Update(req file_proto.ShareRequest) chan FileMsg {
	return lfs.FilesChan
}

// walkDatasourceParents creates helper index, aliases and push the dirs that makes the root path of the datasource
func (lfs *LocalFs) walkDatasourceParents() error {
	// Create index and put mapping if does not exist
	c := db_config_proto.NewConfigClient(globals.DB_SERVICE_NAME, nil)
	_, err := c.CreateIndex(
		context.Background(),
		&db_config_proto.CreateIndexRequest{
			Index: globals.IndexHelper,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.AddAlias(
		context.Background(),
		&db_config_proto.AddAliasRequest{
			Index: globals.IndexHelper,
			Alias: globals.GetMD5Hash(lfs.Endpoint.UserId),
		},
	)
	if err != nil {
		return err
	}

	// Generate files from root to datasource entry point
	pathHelper := strings.Split(lfs.RootPath[1:], "/")
	path := ""

	for i := 0; i < len(pathHelper)-1; i++ {
		path += "/" + pathHelper[i]
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}

		f := file.NewKazoupFileFromLocal(&local.LocalFile{
			Path: path,
			Info: info,
		}, lfs.Endpoint.Id, lfs.Endpoint.UserId, globals.IndexHelper)

		lfs.FilesChan <- NewFileMsg(f, nil)
	}

	return nil
}

// walkHandler recursively discover files belonging to datasource
func (lfs *LocalFs) walkHandler() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print("Got error %s", err)
			return nil
		}
		select {
		case <-lfs.WalkRunning:
			log.Print("Scanner stopped")
			return errors.New("Scanner stopped")
		default:
			f := file.NewKazoupFileFromLocal(&local.LocalFile{
				Path: path,
				Info: info,
			}, lfs.Endpoint.Id, lfs.Endpoint.UserId, lfs.Endpoint.Index)

			lfs.FilesChan <- NewFileMsg(f, nil)
		}

		return nil
	}
}
