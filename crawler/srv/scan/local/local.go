package local

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	scan "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/structs"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Local ...
type Local struct {
	Id       int64
	RootPath string
	Running  chan bool
	Config   map[string]string
	Scanner  scan.Scanner
}

const topic string = "go.micro.topic.files"

// NewLocal ...
func NewLocal(id int64, rootPath string, conf map[string]string) (*Local, error) {
	return &Local{
		Id:       id,
		RootPath: path.Clean(rootPath),
		Running:  make(chan bool, 1),
		Config:   conf,
	}, nil
}

// Start ...
func (fs *Local) Start(crawls map[int64]scan.Scanner, index int64) {
	go func() {
		filepath.Walk(fs.RootPath, fs.walkHandler())
		// Local scan finished
		fs.Stop()
		delete(crawls, index)
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
				return errors.New("Error marshaling data")
			}

			msg := &crawler.FileMessage{
				Id:   f.URL,
				Data: string(b),
			}

			ctx := context.TODO()
			if err := client.Publish(ctx, client.NewPublication(topic, msg)); err != nil {
				log.Printf("Error pubslishing : %", err.Error())
				return err
			}

		}

		return nil
	}
}
