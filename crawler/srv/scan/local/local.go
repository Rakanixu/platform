package local

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	scan "github.com/kazoup/platform/crawler/srv/scan"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Local ...
type Local struct {
	Id       int64
	RootPath string
	Index    string
	Running  chan bool
	Config   map[string]string
	Scanner  scan.Scanner
}

const (
	topic       = "go.micro.topic.files"
	indexHelper = "files_helper"
	filesAlias  = "files"
	fileType    = "file"
)

// NewLocal ...
func NewLocal(id int64, rootPath string, index string, conf map[string]string) (*Local, error) {
	return &Local{
		Id:       id,
		RootPath: path.Clean(rootPath),
		Index:    index,
		Running:  make(chan bool, 1),
		Config:   conf,
	}, nil
}

// Start ...
func (fs *Local) Start(crawls map[int64]scan.Scanner, ds int64) {
	go func() {
		fs.walkDatasourceParents()
		filepath.Walk(fs.RootPath, fs.walkHandler())
		// Local scan finished
		fs.Stop()
		delete(crawls, ds)
	}()
}

// Stop ...
func (fs *Local) Stop() {
	fs.Running <- false
}

// Info ...
func (fs *Local) Info() (scan.Info, error) {
	return scan.Info{
		Id:          fs.Id,
		Type:        "filescanner",
		Description: "File system scanner",
		Config:      fs.Config,
	}, nil
}

func (fs *Local) walkDatasourceParents() error {
	// Create index and put mapping if does not exist
	c := db_proto.NewDBClient("", nil)

	_, err := c.CreateIndexWithSettings(
		context.Background(),
		&db_proto.CreateIndexWithSettingsRequest{
			Index: indexHelper,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.PutMappingFromJSON(
		context.Background(),
		&db_proto.PutMappingFromJSONRequest{
			Index: indexHelper,
			Type:  fileType,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.AddAlias(
		context.Background(),
		&db_proto.AddAliasRequest{
			Index: indexHelper,
			Alias: filesAlias,
		},
	)
	if err != nil {
		return err
	}

	// Generate files from root to datasource entry point
	pathHelper := strings.Split(fs.RootPath[1:], "/")
	path := ""

	for i := 0; i < len(pathHelper)-1; i++ {
		path += "/" + pathHelper[i]
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}

		f := structs.NewDesktopFile(&structs.LocalFile{
			Type: "LocalFile",
			Path: path,
			Info: info,
		})

		b, err := json.Marshal(f)
		if err != nil {
			return err
		}

		msg := &crawler.FileMessage{
			Id:    getMD5Hash(f.URL),
			Index: indexHelper,
			Data:  string(b),
		}

		ctx := context.TODO()
		if err := client.Publish(ctx, client.NewPublication(topic, msg)); err != nil {
			return err
		}
	}

	return nil
}

func (fs *Local) walkHandler() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print("Got error %s", err)
			return nil
		}
		select {
		case <-fs.Running:
			log.Print("Scanner stopped")
			return errors.New("Scanner stopped")
		default:
			f := structs.NewDesktopFile(&structs.LocalFile{
				Type: "LocalFile",
				Path: path,
				Info: info,
			})

			b, err := json.Marshal(f)
			if err != nil {
				return err
			}

			msg := &crawler.FileMessage{
				Id:    getMD5Hash(f.URL),
				Index: fs.Index,
				Data:  string(b),
			}

			ctx := context.TODO()
			if err := client.Publish(ctx, client.NewPublication(topic, msg)); err != nil {
				return err
			}

		}

		return nil
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
