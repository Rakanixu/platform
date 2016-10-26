package engine

import (
	"encoding/json"
	"errors"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

const (
	localEndpoint      = "local://"
	googledriveEnpoint = "googledrive://"
	gmailEnpoint       = "gmail://"
	onedriveEndpoint   = "onedrive://"
	slackEnpoint       = "slack://"
	dropboxEnpoint     = "dropbox://"
	boxEnpoint         = "box://"
)

// Engine interface implements Validation and Save for datasources, name probably could be better DataSourcerer?? jaja
type Engine interface {
	Validate(datasources string) (*datasource_proto.Endpoint, error)
	Save(ctx context.Context, data interface{}, id string) error
	Delete(ctx context.Context, c client.Client) error
}

// NewDataSourceEngine returns a Engine interface
func NewDataSourceEngine(endpoint *datasource_proto.Endpoint) (Engine, error) {
	if strings.Contains(endpoint.Url, localEndpoint) {
		return &Local{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, googledriveEnpoint) {
		return &Googledrive{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, gmailEnpoint) {
		return &Gmail{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, onedriveEndpoint) {
		return &Onedrive{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, slackEnpoint) {
		return &Slack{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, dropboxEnpoint) {
		return &Dropbox{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, boxEnpoint) {
		return &Box{
			Endpoint: *endpoint,
		}, nil
	}

	err := errors.New("Error parsing endpoint for " + endpoint.Url)

	return nil, err
}

// SaveDataSource is a helper to write DS in ES.
func SaveDataSource(ctx context.Context, data interface{}, id string) error {
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

// DeleteDataSource deletes a datasource previously stored and index associated with it
func DeleteDataSource(ctx context.Context, c client.Client, endpoint *datasource_proto.Endpoint) error {
	if endpoint != nil {
		// Delete record from datasources index
		srvReq := c.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: "datasources",
				Type:  "datasource",
				Id:    endpoint.Id,
			},
		)
		srvRes := &db_proto.DeleteResponse{}

		if err := c.Call(ctx, srvReq, srvRes); err != nil {
			return err
		}

		// Remove index for datasource associated with it
		deleteIndexReq := c.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.DeleteIndex",
			&db_proto.DeleteIndexRequest{
				Index: endpoint.Index,
			},
		)
		deleteIndexRes := &db_proto.DeleteResponse{}

		if err := c.Call(ctx, deleteIndexReq, deleteIndexRes); err != nil {
			return err
		}
	}

	return nil
}
