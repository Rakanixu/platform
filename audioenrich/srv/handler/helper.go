package handler

import (
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

func retrieveFilesNotProcessed(ctx context.Context, c client.Client, dsID, index string) ([]string, error) {

	rsp, err := db_helper.SearchFromDB(c, ctx, &db_proto.SearchRequest{

	})
	if err != nil {
		return []string{}, err
	}

	var

	rsp.Result

	return []string{}, nil
}
