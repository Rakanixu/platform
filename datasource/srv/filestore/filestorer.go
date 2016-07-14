package filestorer

import (
	"encoding/json"

	elastic "github.com/kazoup/platform/elastic/srv/proto/elastic"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type FileStorer interface {
	Validate() error
	Save(data interface{}) error
}

type FileStore struct {
	FileStorer
}

// Save FileStore configuration
func (fs *FileStore) Save(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elasticsearch.Create",
		&elastic.CreateRequest{
			Index: "datasources",
			Type:  "datasource",
			Data:  string(b),
		},
	)
	srvRes := &elastic.CreateResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return err
	}

	return nil
}
