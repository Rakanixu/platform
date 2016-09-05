package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/kazoup/platform/datasource/srv/filestore"
	fake "github.com/kazoup/platform/datasource/srv/filestore/fake"
	googledrive "github.com/kazoup/platform/datasource/srv/filestore/googledrive"
	local "github.com/kazoup/platform/datasource/srv/filestore/local"
	onedrive "github.com/kazoup/platform/datasource/srv/filestore/onedrive"
	slack "github.com/kazoup/platform/datasource/srv/filestore/slack"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
	"strings"
)

const (
	fakeEndpoint       = "fake"
	localEndpoint      = "local://"
	googledriveEnpoint = "googledrive://"
	onedriveEndpoint   = "onedrive://"
	slackEnpoint       = "slack://"
	nfsEndpoint        = "nfs://"
	smbEndpoint        = "smb://"
	filesHelperIndex   = "files_helper"
)

// GetDataSource returns a FileStorer interface
func GetDataSource(ds *DataSource, endpoint *proto.Endpoint) (filestorer.FileStorer, error) {
	if strings.Contains(endpoint.Url, fakeEndpoint) {
		return &fake.Fake{
			FileStore: filestorer.FileStore{
				ElasticServiceName: ds.ElasticServiceName,
			},
		}, nil
	}

	if strings.Contains(endpoint.Url, localEndpoint) {
		return &local.Local{
			Endpoint: *endpoint,
			FileStore: filestorer.FileStore{
				ElasticServiceName: ds.ElasticServiceName,
			},
		}, nil
	}

	if strings.Contains(endpoint.Url, googledriveEnpoint) {
		return &googledrive.Googledrive{
			FileStore: filestorer.FileStore{
				ElasticServiceName: ds.ElasticServiceName,
			},
		}, nil
	}
	if strings.Contains(endpoint.Url, onedriveEndpoint) {
		return &onedrive.Onedrive{
			FileStore: filestorer.FileStore{
				ElasticServiceName: ds.ElasticServiceName,
			},
		}, nil
	}
	if strings.Contains(endpoint.Url, slackEnpoint) {
		return &slack.Slack{
			Endpoint: *endpoint,
			FileStore: filestorer.FileStore{
				ElasticServiceName: ds.ElasticServiceName,
			},
		}, nil
	}

	if strings.Contains(endpoint.Url, smbEndpoint) {
		//return &blabla{}, nil
	}

	err := errors.New("Error parsing endpoint for " + endpoint.Url)

	return nil, err
}

// DeleteDataSource deletes a datasource previously stored and index associated with it
func DeleteDataSource(ds *DataSource, id string) error {
	var endpoint *proto.Endpoint

	// Get datasource
	readReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	readRes := &db_proto.ReadResponse{}

	if err := client.Call(context.Background(), readReq, readRes); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(readRes.Result), &endpoint); err != nil {
		return err
	}

	if endpoint != nil {
		// Specific clean up for local datasources
		if strings.Contains(endpoint.Url, localEndpoint) {
			// Remove records from helper index that only belongs to the datasource
			if err := cleanFilesHelperIndex(ds, endpoint); err != nil {
				return err
			}
		}

		// Delete record from datasources index
		srvReq := client.NewRequest(
			ds.ElasticServiceName,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: "datasources",
				Type:  "datasource",
				Id:    id,
			},
		)
		srvRes := &db_proto.DeleteResponse{}

		if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
			return err
		}

		// Remove index for datasource associated with it
		deleteIndexReq := client.NewRequest(
			ds.ElasticServiceName,
			"DB.DeleteIndex",
			&db_proto.DeleteIndexRequest{
				Index: endpoint.Index,
			},
		)
		deleteIndexRes := &db_proto.DeleteResponse{}

		if err := client.Call(context.Background(), deleteIndexReq, deleteIndexRes); err != nil {
			return err
		}
	}

	return nil
}

// SearchDataSources queries for datasources stored in ES
func SearchDataSources(ds *DataSource, req *proto.SearchRequest) (*proto.SearchResponse, error) {
	srvReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.Search",
		&db_proto.SearchRequest{
			Index:    "datasources",
			Type:     "datasource",
			From:     req.From,
			Size:     req.Size,
			Category: req.Category,
			Term:     req.Term,
			Depth:    req.Depth,
			Url:      req.Url,
		},
	)
	srvRes := &db_proto.SearchResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return nil, err
	}

	rsp := &proto.SearchResponse{
		Result: srvRes.Result,
		Info:   srvRes.Info,
	}

	return rsp, nil
}

func ScanDataSource(ds *DataSource, ctx context.Context, id string) error {
	dbSrvReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	dbSrvRes := &db_proto.ReadResponse{}

	if err := client.Call(ctx, dbSrvReq, dbSrvRes); err != nil {
		return err
	}

	var endpoint *proto.Endpoint
	dec := json.NewDecoder(bytes.NewBufferString(dbSrvRes.Result))
	if err := dec.Decode(&endpoint); err != nil {
		return err
	}

	msg := ds.Client.NewPublication(
		globals.ScanTopic,
		endpoint,
	)

	if err := ds.Client.Publish(ctx, msg); err != nil {
		return err
	}

	return nil
}

func CreateIndexWithAlias(ds *DataSource, ctx context.Context, endpoint *proto.Endpoint) error {
	// Create index
	createIndexSrvReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.CreateIndexWithSettings",
		&db_proto.CreateIndexWithSettingsRequest{
			Index: endpoint.Index,
		},
	)
	createIndexSrvRes := &db_proto.CreateIndexWithSettingsResponse{}

	if err := client.Call(ctx, createIndexSrvReq, createIndexSrvRes); err != nil {
		return err
	}

	// Put mapping
	mappingSrvReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.PutMappingFromJSON",
		&db_proto.PutMappingFromJSONRequest{
			Index: endpoint.Index,
			Type:  "file",
		},
	)
	mappingSrvRes := &db_proto.PutMappingFromJSONResponse{}

	if err := client.Call(ctx, mappingSrvReq, mappingSrvRes); err != nil {
		return err
	}

	// Create DS alias
	addAliasReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: endpoint.Id,
		},
	)
	addAliasRes := &db_proto.AddAliasResponse{}

	if err := client.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	// Create specific "files" alias
	addAliasReq = client.NewRequest(
		ds.ElasticServiceName,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: "files",
		},
	)

	if err := client.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	return nil
}

func cleanFilesHelperIndex(ds *DataSource, endpoint *proto.Endpoint) error {
	var datasources []*proto.Endpoint

	rsp, err := SearchDataSources(ds, &proto.SearchRequest{
		Index: "datasources",
		Type:  "datasource",
		From:  0,
		Size:  9999,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rsp.Result), &datasources); err != nil {
		return err
	}

	idx := strings.LastIndex(endpoint.Url, "/")
	if idx > 0 {
		deleteZombieRecords(ds, datasources, endpoint.Url[:idx])
	}

	return nil
}

func deleteZombieRecords(ds *DataSource, datasources []*proto.Endpoint, urlToDelete string) {
	delete := 0

	for _, v := range datasources {
		if !strings.Contains(v.Url, urlToDelete) {
			log.Println("something to delete", v.Url, urlToDelete)
			delete++
		}
	}

	if delete >= len(datasources)-1 {
		deleteReq := client.NewRequest(
			ds.ElasticServiceName,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: filesHelperIndex,
				Type:  "file",
				Id:    getMD5Hash(urlToDelete[len(localEndpoint):]),
			},
		)
		deleteRes := &db_proto.DeleteResponse{}

		if err := client.Call(context.Background(), deleteReq, deleteRes); err != nil {
			log.Println("ERROR", err)
		}
		idx := strings.LastIndex(urlToDelete, "/")

		if idx > 0 && urlToDelete[:idx] != "local:/" {
			deleteZombieRecords(ds, datasources, urlToDelete[:idx])
		}
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
