package local

import (
	"encoding/json"
	"errors"
	scan "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/structs"
	"github.com/micro/go-micro/client"
	example "github.com/micro/micro/examples/template/srv/proto/example"
	"golang.org/x/net/context"
	"log"
	"os"
	"path"
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

	log.Print(path.Clean(rootPath))
	return &Local{
		Id:       id,
		RootPath: path.Clean(rootPath),
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

			//time.Sleep(1 * time.Millisecond)
			msg := &example.Message{
				Say: string(b),
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
