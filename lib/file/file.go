package file

import (
	"encoding/json"
	"github.com/kazoup/platform/crawler/srv/proto/crawler"
	"github.com/kazoup/platform/lib/globals"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"time"
)

const (
	DEFAULT_IMAGE_PREVIEW_URL string = globals.SERVER_ADDRESS + "/media/image/http?source=http://www.scaleautomag.com/sitefiles/images/no-preview-available.png"
	BASE_URL_FILE_PREVIEW     string = globals.SERVER_ADDRESS + "/media"
)

type File interface {
	PreviewURL(width, height, mode, quality string) string
	GetID() string
	GetName() string
	GetUserID() string
	GetIDFromOriginal() string
	GetIndex() string
	GetDatasourceID() string
	GetFileType() string
	GetPathDisplay() string
	GetURL() string
	GetExtension() string
	GetBase64() string
	GetModifiedTime() time.Time
	GetContent() string
	GetOptsTimestamps() *OptsKazoupFile
	SetOptsTimestamps(optsKazoupFile *OptsKazoupFile)
	SetHighlight(highlight string)
	SetContentCategory(contentCategory string)
	SetEntities(entities *rossetelib.RosseteEntities)
}

func IndexAsync(ctx context.Context, c client.Client, file File, topic, index string, notify bool) error {
	b, err := json.Marshal(file)
	if err != nil {
		return err
	}

	msg := &crawler.FileMessage{
		Id:     file.GetID(),
		UserId: file.GetUserID(),
		Index:  index,
		Notify: notify,
		Data:   string(b),
	}

	if err := c.Publish(ctx, c.NewPublication(topic, msg)); err != nil {
		return err
	}

	return nil
}
