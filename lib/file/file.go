package file

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

const (
	DEFAULT_IMAGE_PREVIEW_URL string = globals.SERVER_ADDRESS + "/media/image/http?source=http://www.scaleautomag.com/sitefiles/images/no-preview-available.png"
	BASE_URL_FILE_PREVIEW     string = globals.SERVER_ADDRESS + "/media"
)

type File interface {
	PreviewURL(width, height, mode, quality string) string
	GetID() string
	GetIDFromOriginal() string
	GetIndex() string
	GetDatasourceID() string
	GetFileType() string
	GetPathDisplay() string
	GetURL() string
}

func IndexAsync(file File, topic, index string) error {
	b, err := json.Marshal(file)
	if err != nil {
		return err
	}

	msg := &crawler.FileMessage{
		Id:    file.GetID(),
		Index: index,
		Data:  string(b),
	}

	if err := client.Publish(context.Background(), client.NewPublication(topic, msg)); err != nil {
		return err
	}

	return nil
}
