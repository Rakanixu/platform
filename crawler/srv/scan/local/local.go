package local

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/go-homedir"
	scan "github.com/kazoup/platform/crawler/srv/scan"
	publish "github.com/kazoup/platform/publish/srv/proto/publish"
	"github.com/kazoup/platform/structs"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"os"
	"path/filepath"
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
	path, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	return &Local{
		Id:       id,
		RootPath: path,
		Running:  make(chan bool, 1),
		Config:   conf,
	}, nil
}

// Start ...
func (fs *Local) Start() {
	go filepath.Walk(fs.RootPath, fs.walkHandler())
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
		select {
		case <-fs.Running:
			return errors.New("Scanner stopped")
		default:
			f := structs.NewFileFromLocal(&structs.LocalFile{
				Type: "LocalFile",
				Path: "/" + path,
				Info: info,
			})

			b, err := json.Marshal(f)
			if err != nil {
				return errors.New("Error marshaling data")
			}

			req := client.NewRequest(
				"go.micro.srv.publish",
				"Publish.Send",
				&publish.SendRequest{
					Topic: topic,
					Data:  string(b),
				},
			)
			res := &publish.SendResponse{}

			// Call Publish.Send
			if err := client.Call(context.Background(), req, res); err != nil {
				return errors.New("Error calling com.kazoup.srv.publish.Publish.Send")
			}
		}

		return nil
	}
}
