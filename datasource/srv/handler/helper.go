package handler

import (
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/filestore"
	local "github.com/kazoup/platform/datasource/srv/filestore/local"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	elastic "github.com/kazoup/platform/elastic/srv/proto/elastic"
	publish "github.com/kazoup/platform/publish/srv/proto/publish"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"strings"
)

const (
	localEndpoint = "local://"
	nfsEndpoint   = "nfs://"
	smbEndpoint   = "smb://"
	topic         = "go.micro.topic.scan"
)

// GetDataSource returns a FileStorer interface
func GetDataSource(endpoint *proto.Endpoint) (filestorer.FileStorer, error) {
	if strings.Contains(endpoint.Url, localEndpoint) {
		return &local.Local{
			Endpoint: *endpoint,
		}, nil
	}

	if strings.Contains(endpoint.Url, nfsEndpoint) {
		//return &blabla{}, nil
	}

	if strings.Contains(endpoint.Url, smbEndpoint) {
		//return &blabla{}, nil
	}

	err := errors.New("Error parsing endpoint for " + endpoint.Url)

	return nil, err
}

// DeleteDataSource deletes a datasource previously stored
func DeleteDataSource(id string) error {
	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Delete",
		&elastic.DeleteRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	srvRes := &elastic.CreateResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return err
	}

	return nil
}

// SearchDataSources queries for datasources stored in ES
func SearchDataSources(query string, limit int64, offset int64) ([]*proto.Endpoint, error) {
	var result []*proto.Endpoint
	var resultMap map[string]interface{}

	if len(query) <= 0 {
		query = "*"
	}

	srvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Search",
		&elastic.SearchRequest{
			Index:  "datasources",
			Type:   "datasource",
			Query:  query,
			Limit:  limit,
			Offset: offset,
		},
	)
	srvRes := &elastic.SearchResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(srvRes.Result), &resultMap); err != nil {
		return nil, err
	}

	var size, files, lastScan int64
	for _, v := range resultMap["hits"].(map[string]interface{})["hits"].([]interface{}) {
		// Sanitize possible empty values
		if v.(map[string]interface{})["_source"].(map[string]interface{})["size"] == nil {
			size = 0
		} else {
			size = int64(v.(map[string]interface{})["_source"].(map[string]interface{})["size"].(float64))
		}
		if v.(map[string]interface{})["_source"].(map[string]interface{})["files"] == nil {
			files = 0
		} else {
			files = int64(v.(map[string]interface{})["_source"].(map[string]interface{})["files"].(float64))
		}
		if v.(map[string]interface{})["_source"].(map[string]interface{})["last_scan"] == nil {
			lastScan = 0
		} else {
			lastScan = int64(v.(map[string]interface{})["_source"].(map[string]interface{})["last_scan"].(float64))
		}

		result = append(result, &proto.Endpoint{
			Id:       v.(map[string]interface{})["_id"].(string),
			Url:      v.(map[string]interface{})["_source"].(map[string]interface{})["url"].(string),
			Size:     size,
			Files:    files,
			LastScan: lastScan,
		})
	}

	return result, nil
}

func ScanDataSource(ctx context.Context, id string) error {
	elasticSrvReq := client.NewRequest(
		"go.micro.srv.elastic",
		"Elastic.Read",
		&elastic.ReadRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	elasticSrvRes := &elastic.ReadResponse{}

	if err := client.Call(ctx, elasticSrvReq, elasticSrvRes); err != nil {
		return err
	}

	pubSrvReq := client.NewRequest(
		"go.micro.srv.publish",
		"Publish.Send",
		&publish.SendRequest{
			Topic: topic,
			Data:  elasticSrvRes.Result,
		},
	)
	pubSrvRes := &publish.SendResponse{}

	if err := client.Call(ctx, pubSrvReq, pubSrvRes); err != nil {
		return err
	}

	return nil
}
