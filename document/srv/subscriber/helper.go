package subscriber

import (
	"encoding/json"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/db/operations"
	"github.com/kazoup/platform/lib/db/operations/proto/operations"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

type info struct {
	Total int64 `json:"total"`
}

var from, size int64

func publishDocFilesNotProcessed(ctx context.Context, endpoint *proto_datasource.Endpoint) error {
	srv, ok := micro.FromContext(ctx)
	if !ok {
		return errors.ErrInvalidCtx
	}

	ids, err := retrieveDocFilesNotProcessed(ctx, srv.Client(), endpoint.Index)
	if err != nil {
		return err
	}

	// Publish msg for all files not being process yet
	for _, v := range ids {
		if err := srv.Client().Publish(ctx, srv.Client().NewPublication(globals.DocEnrichTopic, &enrich.EnrichMessage{
			Index:  endpoint.Index,
			Id:     v,
			UserId: endpoint.UserId,
		})); err != nil {
			log.Println("ERROR publishing DocEnrichTopic message", err)
		}
	}

	return nil
}

func retrieveDocFilesNotProcessed(ctx context.Context, c client.Client, index string) ([]string, error) {
	var ids []string
	from = 0
	size = 9999

	result, err := retrieveFiles(ctx, c, ids, index, from, size)
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

func retrieveFiles(ctx context.Context, c client.Client, ids []string, index string, from, size int64) ([]string, error) {
	rsp, err := operations.Search(ctx, &proto_operations.SearchRequest{
		Index:                index,
		From:                 from,
		Size:                 size,
		Category:             globals.CATEGORY_DOCUMENT,
		Type:                 globals.FileType,
		NoKazoupFileOriginal: true,
	})
	if err != nil {
		return []string{}, err
	}

	var i info
	if err := json.Unmarshal([]byte(rsp.Info), &i); err != nil {
		return []string{}, err
	}

	var r []*file.KazoupFile
	if err := json.Unmarshal([]byte(rsp.Result), &r); err != nil {
		return []string{}, err
	}

	for _, v := range r {
		if v.OptsKazoupFile == nil || v.OptsKazoupFile.ContentTimestamp == nil || v.OptsKazoupFile.ContentTimestamp.Before(v.Modified) {
			ids = append(ids, v.ID)
		}
	}

	if from+size < i.Total {
		ids, err = retrieveFiles(ctx, c, ids, index, from+size, size)
		if err != nil {
			return []string{}, err
		}
	}
	return ids, nil
}
