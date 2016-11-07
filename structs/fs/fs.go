package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

// Fs File System interface. Fyle system is responsible to manage its own files
type Fs interface {
	List() (chan file.File, chan bool, error)
	CreateFile(file_proto.CreateRequest) (*file_proto.CreateResponse, error)
	DeleteFile(context.Context, client.Client, file_proto.DeleteRequest) (*file_proto.DeleteResponse, error)
	ShareFile(context.Context, client.Client, file_proto.ShareRequest) (string, error)
	GetDatasourceId() string
	Token() string
	GetThumbnail(id string) (string, error)
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

	return nil, errors.New("Error parsing URL")
}
