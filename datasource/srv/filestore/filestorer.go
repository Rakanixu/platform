package filestorer

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type FileStorer interface {
	Validate(datasources string) (*datasource_proto.Endpoint, error)
	Save(data interface{}, id string) error
}

type FileStore struct {
	FileStorer         FileStorer
	ElasticServiceName string
}

// Save FileStore configuration
func (fs *FileStore) Save(data interface{}, id string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	srvReq := client.NewRequest(
		fs.ElasticServiceName,
		"DB.Create",
		&db_proto.CreateRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
			Data:  string(b),
		},
	)
	srvRes := &db_proto.CreateResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return err
	}

	return nil
}
