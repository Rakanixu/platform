package handler

import (
	"encoding/json"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	db_helper "github.com/kazoup/platform/lib/dbhelper"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"log"
)

func retrieveAudioFilesNotProcessed(ctx context.Context, c client.Client, dsID, index string) ([]string, error) {
	rsp, err := db_helper.SearchFromDB(c, ctx, &db_proto.SearchRequest{
		Index:                index,
		From:                 0,
		Size:                 1000,
		Category:             globals.CATEGORY_AUDIO,
		Type:                 globals.FileType,
		NoKazoupFileOriginal: true,
	})
	if err != nil {
		return []string{}, err
	}

	var r []*file.KazoupFile

	if err := json.Unmarshal([]byte(rsp.Result), &r); err != nil {
		return []string{}, err
	}

	for _, v := range r {
		log.Println("TIMESTAMP", v.OptsKazoupFile.AudioTimestamp)
	}

	return []string{}, nil
}
