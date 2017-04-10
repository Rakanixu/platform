package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	file "github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"strings"
)

// Fs File System interface. Fyle system is responsible to manage its own files
type Fs interface {
	FsOperations
	FsUtils
	FsProcessing
}

type FsOperations interface {
	Walk() (chan FileMsg, chan bool)
	WalkUsers() (chan UserMsg, chan bool)
	WalkChannels() (chan ChannelMsg, chan bool)
	Create(file_proto.CreateRequest) chan FileMsg
	Delete(file_proto.DeleteRequest) chan FileMsg
	Update(file_proto.ShareRequest) chan FileMsg
}

type FsUtils interface {
	Authorize() (*datasource_proto.Token, error)
	GetDatasourceId() string
	GetThumbnail(string) (string, error)
}

type FsProcessing interface {
	DocEnrich(file.File) chan FileMsg
	ImgEnrich(file.File) chan FileMsg
	AudioEnrich(file.File, *gcslib.GoogleCloudStorage) chan FileMsg
	Thumbnail(file.File, *gcslib.GoogleCloudStorage) chan FileMsg
}

// NewFsFromEndpoint constructor from endpoint
func NewFsFromEndpoint(e *datasource_proto.Endpoint) (Fs, error) {
	dsUrl := strings.Split(e.Url, ":")

	switch dsUrl[0] {
	case globals.Local:
		return NewLocalFsFromEndpoint(e), nil
	case globals.Slack:
		return NewSlackFsFromEndpoint(e), nil
	case globals.GoogleDrive:
		return NewGoogleDriveFsFromEndpoint(e), nil
	case globals.Gmail:
		return NewGmailFsFromEndpoint(e), nil
	case globals.OneDrive:
		return NewOneDriveFsFromEndpoint(e), nil
	case globals.Dropbox:
		return NewDropboxFsFromEndpoint(e), nil
	case globals.Box:
		return NewBoxFsFromEndpoint(e), nil
	default:
		return nil, errors.New("Not such file system (fs)")
	}

}
