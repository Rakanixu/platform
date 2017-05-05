package fs

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	file "github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
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
	Create(proto_file.CreateRequest) chan FileMsg
	Delete(proto_file.DeleteRequest) chan FileMsg
	Update(proto_file.ShareRequest) chan FileMsg
}

type FsUtils interface {
	Authorize() (*proto_datasource.Token, error)
	GetDatasourceId() string
	GetThumbnail(string) (string, error)
}

type FsProcessing interface {
	DocEnrich(file.File) chan FileMsg
	ImgEnrich(file.File) chan FileMsg
	AudioEnrich(file.File) chan FileMsg
	Thumbnail(file.File) chan FileMsg
}

// NewFsFromEndpoint constructor from endpoint
func NewFsFromEndpoint(e *proto_datasource.Endpoint) (Fs, error) {
	dsUrl := strings.Split(e.Url, ":")

	switch dsUrl[0] {
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

// UpdateFsAuth is a helper for updating auth details
func UpdateFsAuth(ctx context.Context, id string, token *proto_datasource.Token) error {
	var ds *proto_datasource.Endpoint

	// Retrieve endpoint
	rr, err := operations.Read(ctx, &proto_operations.ReadRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    id,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rr.Result), &ds); err != nil {
		return err
	}

	ds.Token = token

	b, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	_, err = operations.Update(ctx, &proto_operations.UpdateRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}

// ClearIndex is a helper that remove records (Files) from db that not longer belong to a datasource
// Compares LastSeen with the time the crawler started
// so all records with a LastSeen before will be removed from index
// file does not exists any more on datasource
func ClearIndex(ctx context.Context, e *proto_datasource.Endpoint) error {
	_, err := operations.DeleteByQuery(ctx, &proto_operations.DeleteByQueryRequest{
		Indexes:  []string{e.Index},
		Types:    []string{globals.FileType},
		LastSeen: e.LastScanStarted,
	})
	if err != nil {
		return err
	}

	return nil
}
