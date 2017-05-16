package engine

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/config"
	"github.com/kazoup/platform/lib/db/config/proto/config"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/utils"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"strings"
)

const (
	googledriveEnpoint = "googledrive://"
	gmailEnpoint       = "gmail://"
	onedriveEndpoint   = "onedrive://"
	slackEnpoint       = "slack://"
	dropboxEnpoint     = "dropbox://"
	boxEnpoint         = "box://"
)

// Engine interface implements Validation, Save, Delete, Update for datasources
type Engine interface {
	Validate(ctx context.Context, datasources string) (*proto_datasource.Endpoint, error)
	Save(ctx context.Context, data interface{}, id string) error
	Delete(ctx context.Context) error
	CreateIndexWithAlias(ctx context.Context) error
}

// NewDataSourceEngine returns a Engine interface
func NewDataSourceEngine(endpoint *proto_datasource.Endpoint) (Engine, error) {
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

	if strings.Contains(endpoint.Url, globals.Mock) {
		return &Mock{
			Endpoint: *endpoint,
		}, nil
	}

	err := errors.New("com.kazoup.Datasource.NewDataSourceEngine", "Error parsing endpoint for ", 500)

	return nil, err
}

// GenerateEndpoint assings index and id if data does not exists
func GenerateEndpoint(ctx context.Context, endpoint proto_datasource.Endpoint) (proto_datasource.Endpoint, error) {
	// Search for existing datasource. If exists, use instance.
	// DB.Search internally tries to get userId explicitly (from request). If not preset, tries to get from context
	// If context is system context, an error will be trow
	srvRsp, err := SearchDataSources(ctx, &proto_datasource.SearchRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		From:  0,
		Size:  9999,
		Url:   endpoint.Url,
	})
	if err != nil {
		return proto_datasource.Endpoint{}, err
	}

	var r []*proto_datasource.Endpoint
	if err := json.Unmarshal([]byte(srvRsp.Result), &r); err != nil {
		return proto_datasource.Endpoint{}, err
	}
	// user_id + url makes a constraint, so record should be unique.
	if len(r) == 1 {
		endpoint = *r[0]
	}

	if len(endpoint.Index) == 0 {
		str, err := utils.NewUUID()
		if err != nil {
			return endpoint, err
		}

		endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	}

	if len(endpoint.Id) == 0 {
		endpoint.Id = utils.GetMD5Hash(endpoint.Url + endpoint.UserId)
	}

	return endpoint, nil
}

// SaveDataSource is a helper to write DS in ES.
func SaveDataSource(ctx context.Context, data interface{}, id string) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = operations.Create(ctx, &proto_operations.CreateRequest{
		Index: globals.IndexDatasources,
		Type:  globals.TypeDatasource,
		Id:    id,
		Data:  string(b),
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteDataSource deletes a datasource previously stored and index associated with it
func DeleteDataSource(ctx context.Context, endpoint *proto_datasource.Endpoint) error {
	if endpoint != nil {
		// Remove index for datasource associated with it
		_, err := config.DeleteIndex(ctx, &proto_config.DeleteIndexRequest{
			Index: endpoint.Index,
		})
		if err != nil {
			return err
		}

		// Delete record from datasources index
		_, err = operations.Delete(ctx, &proto_operations.DeleteRequest{
			Index: globals.IndexDatasources,
			Type:  globals.TypeDatasource,
			Id:    endpoint.Id,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func CreateIndexWithAlias(ctx context.Context, endpoint *proto_datasource.Endpoint) error {
	// Create index
	_, err := config.CreateIndex(ctx, &proto_config.CreateIndexRequest{
		Index: endpoint.Index,
	})
	if err != nil {
		return err
	}

	// Create DS alias
	_, err = config.AddAlias(ctx, &proto_config.AddAliasRequest{
		Index: endpoint.Index,
		Alias: endpoint.Id,
	})
	if err != nil {
		return err
	}

	// Create specific "files" alias
	_, err = config.AddAlias(ctx, &proto_config.AddAliasRequest{
		Index: endpoint.Index,
		Alias: utils.GetMD5Hash(endpoint.UserId),
	})
	if err != nil {
		return err
	}

	return nil
}
