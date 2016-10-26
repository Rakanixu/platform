package handler

import (
	proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	"golang.org/x/net/context"
)

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
