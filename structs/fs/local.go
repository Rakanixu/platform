package fs

import (
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	file_proto "github.com/kazoup/platform/file/srv/proto/file"
	"github.com/kazoup/platform/structs/file"
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/local"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LocalFs struct
type LocalFs struct {
	Endpoint  *datasource_proto.Endpoint
	Running   chan bool
	RootPath  string
	FilesChan chan file.File
}

// NewLocalFsFromEndpoint constructor
func NewLocalFsFromEndpoint(e *datasource_proto.Endpoint) Fs {
	url := strings.Split(e.Url, "://")

	return &LocalFs{
		Endpoint:  e,
		Running:   make(chan bool, 1),
		RootPath:  url[1],
		FilesChan: make(chan file.File),
	}
}

// List returns 2 channels, for files and state. Discover local files
func (lfs *LocalFs) List() (chan file.File, chan bool, error) {
	go func() {
		if err := lfs.walkDatasourceParents(); err != nil {
			log.Println("ERROR", err)
		}

		if err := filepath.Walk(lfs.RootPath, lfs.walkHandler()); err != nil {
			log.Println("ERROR", err)
		}
		lfs.Running <- false
	}()

	return lfs.FilesChan, lfs.Running, nil
}

// Token belongs to Fs interface
func (lfs *LocalFs) Token() string {
	// LocalFs cannot have Token, cause represents a Local datasource which does not required oauth
	return ""
}

// GetDatasourceId returns datasource ID
func (lfs *LocalFs) GetDatasourceId() string {
	return lfs.Endpoint.Id
}

// GetThumbnail belongs to Fs interface
func (lfs *LocalFs) GetThumbnail(id string) (string, error) {
	return "", nil
}

// CreateFile belongs to Fs interface
func (lfs *LocalFs) CreateFile(rq file_proto.CreateRequest) (*file_proto.CreateResponse, error) {
	return &file_proto.CreateResponse{}, nil
}

// DeleteFile deletes a local file
func (lfs *LocalFs) DeleteFile(ctx context.Context, c client.Client, rq file_proto.DeleteRequest) (*file_proto.DeleteResponse, error) {
	return &file_proto.DeleteResponse{}, nil
}

// ShareFile
func (lfs *LocalFs) ShareFile(ctx context.Context, c client.Client, req file_proto.ShareRequest) (string, error) {
	return "", nil
}

// walkDatasourceParents creates helper index, aliases and push the dirs that makes the root path of the datasource
func (lfs *LocalFs) walkDatasourceParents() error {
	// Create index and put mapping if does not exist
	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, nil)
	_, err := c.CreateIndexWithSettings(
		context.Background(),
		&db_proto.CreateIndexWithSettingsRequest{
			Index: globals.IndexHelper,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.PutMappingFromJSON(
		context.Background(),
		&db_proto.PutMappingFromJSONRequest{
			Index: globals.IndexHelper,
			Type:  globals.FileType,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.AddAlias(
		context.Background(),
		&db_proto.AddAliasRequest{
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

		lfs.FilesChan <- f
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
		case <-lfs.Running:
			log.Print("Scanner stopped")
			return errors.New("Scanner stopped")
		default:
			f := file.NewKazoupFileFromLocal(&local.LocalFile{
				Path: path,
				Info: info,
			}, lfs.Endpoint.Id, lfs.Endpoint.UserId, lfs.Endpoint.Index)

			lfs.FilesChan <- f
		}

		return nil
	}
}
