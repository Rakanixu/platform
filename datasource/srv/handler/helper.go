package handler

import (
	"bytes"
	"encoding/json"
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
	"log"
	"time"
)

const (
	localEndpoint    = "local://"
	filesHelperIndex = "files_helper"
)

/*// DeleteDataSource deletes a datasource previously stored and index associated with it
func DeleteDataSource(ctx context.Context, ds *DataSource, id string) error {
	var endpoint *proto.Endpoint

	// Get datasource
	readReq := ds.Client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.Read",
		&db_proto.ReadRequest{
			Index: "datasources",
			Type:  "datasource",
			Id:    id,
		},
	)
	readRes := &db_proto.ReadResponse{}

	if err := ds.Client.Call(ctx, readReq, readRes); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(readRes.Result), &endpoint); err != nil {
		return err
	}

	if endpoint != nil {
		// Specific clean up for local datasources
		if strings.Contains(endpoint.Url, localEndpoint) {
			// Remove records from helper index that only belongs to the datasource
			if err := cleanFilesHelperIndex(ctx, ds, endpoint); err != nil {
				return err
			}
		}

		// Delete record from datasources index
		srvReq := ds.Client.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: "datasources",
				Type:  "datasource",
				Id:    id,
			},
		)
		srvRes := &db_proto.DeleteResponse{}

		if err := ds.Client.Call(ctx, srvReq, srvRes); err != nil {
			return err
		}

		// Remove index for datasource associated with it
		deleteIndexReq := ds.Client.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.DeleteIndex",
			&db_proto.DeleteIndexRequest{
				Index: endpoint.Index,
			},
		)
		deleteIndexRes := &db_proto.DeleteResponse{}

		if err := ds.Client.Call(ctx, deleteIndexReq, deleteIndexRes); err != nil {
			return err
		}
	}

	return nil
}*/

/*// SearchDataSources queries for datasources stored in ES
func SearchDataSources(ctx context.Context, ds *DataSource, req *proto.SearchRequest) (*proto.SearchResponse, error) {
	srvReq := ds.Client.NewRequest(
		globals.DB_SERVICE_NAME,
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

	if err := ds.Client.Call(ctx, srvReq, srvRes); err != nil {
		return nil, err
	}

	rsp := &proto.SearchResponse{
		Result: srvRes.Result,
		Info:   srvRes.Info,
	}

	return rsp, nil
}*/

func ScanDataSource(ds *DataSource, ctx context.Context, id string) error {
	//FIXME use ds.Client
	c := db_proto.NewDBClient(globals.DB_SERVICE_NAME, nil)
	dbSrvRes, err := c.Read(ctx, &db_proto.ReadRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    id,
	})
	if err != nil {
		log.Println(err)
	}

	var endpoint *proto.Endpoint
	dec := json.NewDecoder(bytes.NewBufferString(dbSrvRes.Result))
	if err := dec.Decode(&endpoint); err != nil {
		return err
	}

	// Set time for starting scan, crawler running  and update datasource
	endpoint.CrawlerRunning = true
	endpoint.LastScanStarted = time.Now().Unix()
	b, err := json.Marshal(endpoint)
	if err != nil {
		log.Println(err)
	}
	_, err = c.Update(ctx, &db_proto.UpdateRequest{
		Index: "datasources",
		Type:  "datasource",
		Id:    endpoint.Id,
		Data:  string(b),
	})
	if err != nil {
		log.Println(err)
	}

	// Publish scan topic, crawlers should pick up message and start scanning
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
	createIndexSrvReq := ds.Client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.CreateIndexWithSettings",
		&db_proto.CreateIndexWithSettingsRequest{
			Index: endpoint.Index,
		},
	)
	createIndexSrvRes := &db_proto.CreateIndexWithSettingsResponse{}

	if err := ds.Client.Call(ctx, createIndexSrvReq, createIndexSrvRes); err != nil {
		return err
	}

	// Create DS alias
	addAliasReq := ds.Client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: endpoint.Id,
		},
	)
	addAliasRes := &db_proto.AddAliasResponse{}

	if err := ds.Client.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	// Create specific "files" alias
	addAliasReq = ds.Client.NewRequest(
		globals.DB_SERVICE_NAME,
		"DB.AddAlias",
		&db_proto.AddAliasRequest{
			Index: endpoint.Index,
			Alias: globals.GetMD5Hash(endpoint.UserId),
		},
	)

	if err := ds.Client.Call(ctx, addAliasReq, addAliasRes); err != nil {
		return err
	}

	return nil
}

/*
func cleanFilesHelperIndex(ctx context.Context, ds *DataSource, endpoint *proto.Endpoint) error {
	var datasources []*proto.Endpoint

	rsp, err := SearchDataSources(ctx, ds, &proto.SearchRequest{
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
		deleteZombieRecords(ctx, ds, datasources, endpoint.Url[:idx])
	}

	return nil
}*/

/*func deleteZombieRecords(ctx context.Context, ds *DataSource, datasources []*proto.Endpoint, urlToDelete string) {
	delete := 0

	for _, v := range datasources {
		if !strings.Contains(v.Url, urlToDelete) {
			log.Println("something to delete", v.Url, urlToDelete)
			delete++
		}
	}

	if delete >= len(datasources)-1 {
		deleteReq := ds.Client.NewRequest(
			globals.DB_SERVICE_NAME,
			"DB.Delete",
			&db_proto.DeleteRequest{
				Index: filesHelperIndex,
				Type:  "file",
				Id:    globals.GetMD5Hash(urlToDelete[len(localEndpoint):]),
			},
		)
		deleteRes := &db_proto.DeleteResponse{}

		if err := ds.Client.Call(ctx, deleteReq, deleteRes); err != nil {
			log.Println("ERROR", err)
		}
		idx := strings.LastIndex(urlToDelete, "/")

		if idx > 0 && urlToDelete[:idx] != "local:/" {
			deleteZombieRecords(ctx, ds, datasources, urlToDelete[:idx])
		}
	}
}*/
