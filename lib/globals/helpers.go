package globals

import (
	"encoding/json"
	"fmt"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/micro/go-micro/client"
)

func deleteFilesNoExistsFromGCS(e *datasource_proto.Endpoint, from, size int) error {
	// Helper to get only the ids from original file, this way we do  not know what type are we working with
	var r []struct {
		Original struct {
			Id string `json:"id"`
		} `json:"original"`
	}
	var i struct {
		Total int `json:"total"`
	}

	c := db_proto.NewDBClient(DB_SERVICE_NAME, nil)
	rsp, err := c.Search(NewSystemContext(), &db_proto.SearchRequest{
		Index:    e.Index,
		From:     int64(from),
		Size:     int64(size),
		Type:     FileType,
		FileType: FileTypeFile,
		UserId:   e.UserId,
		Category: CATEGORY_PICTURE,
		LastSeen: e.LastScanStarted,
	})
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rsp.Info), &i); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(rsp.Result), &r); err != nil {
		return err
	}

	for _, v := range r {
		// Publish message to clean async the bucket that stores the thumbnails in GC storage
		if err := client.Publish(NewSystemContext(), client.NewPublication(DeleteFileInBucketTopic, &datasource_proto.DeleteFileInBucketMessage{
			FileId: v.Original.Id,
			Index:  e.Index,
		})); err != nil {
			fmt.Println("ERROR cleaningthumbs from GCS", err)
		}
	}

	if i.Total > from+size {
		if err := deleteFilesNoExistsFromGCS(e, from+size, size); err != nil {
			return err
		}
	}

	return nil
}
