package fake

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	scan "github.com/kazoup/platform/crawler/srv/scan"
	"github.com/kazoup/platform/structs"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"time"
)

// Fake ...
type Fake struct {
	Id      int64
	Type    string
	Config  map[string]string
	Running chan bool
	Scanner scan.Scanner
}

// NewFake creates a Fake instance
func NewFake(id int64, conf map[string]string) scan.Scanner {
	return &Fake{
		Id:      id,
		Type:    globals.Fake,
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
					Id:   mockFile.ID,
					Data: string(b),
				}

				ctx := context.TODO()
				if err := client.Publish(ctx, client.NewPublication(globals.FilesTopic, msg)); err != nil {
					log.Printf("Error pubslishing : %", err.Error())
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
