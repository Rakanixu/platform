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
	"time"
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
	Scan(ctx context.Context, c client.Client) error
	CreateIndexWithAlias(ctx context.Context, c client.Client) error
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

// GenerateEndpoint assings index and id if data does not exists
func GenerateEndpoint(endpoint *datasource_proto.Endpoint) (*datasource_proto.Endpoint, error) {
	if len(endpoint.Index) == 0 {
		str, err := globals.NewUUID()
		if err != nil {
			return endpoint, err
		}

		endpoint.Index = "index" + strings.Replace(str, "-", "", 1)
	}

	if len(endpoint.Id) == 0 {
		endpoint.Id = globals.GetMD5Hash(endpoint.Url + endpoint.UserId)
	}

	return endpoint, nil
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

// ScanDataSource is a helper to kick off scans
func ScanDataSource(ctx context.Context, c client.Client, endpoint *datasource_proto.Endpoint) error {
	// Set time for starting scan, crawler running  and update datasource
	endpoint.CrawlerRunning = true
	endpoint.LastScanStarted = time.Now().Unix()
	b, err := json.Marshal(endpoint)
	if err != nil {
		return err
	}

	srvReq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Update",
		&db_proto.UpdateRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    endpoint.Id,
			Data:  string(b),
		},
	)
	srvRes := &db_proto.UpdateResponse{}

	if err := c.Call(ctx, srvReq, srvRes); err != nil {
		return err
	}

	// Publish scan topic, crawlers should pick up message and start scanning
	msg := c.NewPublication(
		globals.ScanTopic,
		endpoint,
	)

	if err := c.Publish(ctx, msg); err != nil {
		return err
	}

	return nil
}

func CreateIndexWithAlias(ctx context.Context, c client.Client, endpoint *datasource_proto.Endpoint) error {
	// Create index
	createIndexSrvReq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.CreateIndexWithSettings",
		&db_proto.CreateIndexWithSettingsRequest{
			Index: endpoint.Index,
		},
	)
	createIndexSrvRes := &db_proto.CreateIndexWithSettingsResponse{}

	if err := c.Call(ctx, createIndexSrvReq, createIndexSrvRes); err != nil {
		return err
	}

	// Create DS alias
	addAliasReq := c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: endpoint.Id,
		},
	)
	addAliasRes := &db_proto.AddAliasResponse{}

	if err := c.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	// Create specific "files" alias
	addAliasReq = c.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: globals.GetMD5Hash(endpoint.UserId),
		},
	)

	if err := c.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	return nil
}
