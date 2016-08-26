package fake

import (
	"encoding/json"
	"fmt"
	"time"

	scan "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/structs"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

// Fake ...
type Fake struct {
	Id      int64
	Type    string
	Config  map[string]string
	Running chan bool
	Scanner scan.Scanner
}

// TODO: topic
const topic string = "go.micro.topic.files"

// NewFake creates a Fake instance
func NewFake(id int64, conf map[string]string) scan.Scanner {
	return &Fake{
		Id:      id,
		Type:    "fake",
		Config:  conf,
		Running: make(chan bool, 1),
	}
}

// Start fake scan
func (f *Fake) Start(crawls map[int64]scan.Scanner, id int64) {
	go func() {
		for {
			select {
			// We will stop execution when we receive a value from channel
			case <-f.Running:
				return
			default:
				mockFile := structs.NewMockFile()
				//ID := mockFile.ID
				b, err := json.Marshal(mockFile)
				if err != nil {
					fmt.Errorf("Error marshaling data")
				}

				//Publish

				msg := &crawler.FileMessage{
					Id:   f.URL,
					Data: string(b),
				}

				ctx := context.TODO()
				if err := client.Publish(ctx, client.NewPublication(topic, msg)); err != nil {
					log.Printf("Error pubslishing : %", err.Error())
					return err
				}
				time.Sleep(time.Second)
			}
		}

	}()
}

// Stop fake scan
func (f *Fake) Stop() {
	f.Running <- false
}

// Info returns fake scanner info
func (f *Fake) Info() (scan.Info, error) {
	return scan.Info{
		Id:          f.Id,
		Type:        f.Type,
		Description: "Fake scanner",
		Config:      f.Config,
	}, nil
}
