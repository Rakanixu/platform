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
	topic              = "go.micro.topic.scan"
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

// DeleteDataSource deletes a datasource previously stored
func DeleteDataSource(ds *DataSource, id string) error {
	srvReq := client.NewRequest(
		ds.ElasticServiceName,
		"DB.Delete",
		&db_proto.DeleteRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	srvRes := &db_proto.CreateResponse{}

	if err := client.Call(context.Background(), srvReq, srvRes); err != nil {
		return err
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

	log.Println("ss", dbSrvRes)
	log.Println("ss", dbSrvRes.Result)

	var endpoint *proto.Endpoint

	//bytes.NewBufferString("your string")
	dec := json.NewDecoder(bytes.NewBufferString(dbSrvRes.Result))
	if err := dec.Decode(&endpoint); err != nil {
		return err
	}

	log.Println(endpoint)
	/*	if err := json.Unmarshal([]byte(dbSrvRes.Result), &endpoint); err != nil {
		return err
	}*/

	msg := ds.Client.NewPublication(
		topic,
		endpoint,
	)

	if err := ds.Client.Publish(ctx, msg); err != nil {
		return err
	}

	return nil
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
