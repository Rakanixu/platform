package handler

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type info struct {
	Total int64 `json:"total"`
}

var from, size int64

func retrieveAudioFilesNotProcessed(ctx context.Context, c client.Client, dsID, index string) ([]string, error) {
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
	rsp, err := db_helper.SearchFromDB(c, ctx, &db_proto.SearchRequest{
		Index:                index,
		From:                 from,
		Size:                 size,
		Category:             globals.CATEGORY_AUDIO,
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
		if v.OptsKazoupFile == nil || v.OptsKazoupFile.AudioTimestamp == nil || v.OptsKazoupFile.AudioTimestamp.Before(v.Modified) {
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
