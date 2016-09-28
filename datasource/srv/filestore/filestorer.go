package filestorer

import (
	"encoding/json"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type FileStorer interface {
	Validate(datasources string) (*datasource_proto.Endpoint, error)
	Save(ctx context.Context, data interface{}, id string) error
}

type FileStore struct {
	FileStorer FileStorer
}

// Save FileStore configuration
func (fs *FileStore) Save(ctx context.Context, data interface{}, id string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	srvReq := client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Create",
		&db_proto.CreateRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
			Data:  string(b),
		},
	)
	srvRes := &db_proto.CreateResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return err
	}

	return nil
}
