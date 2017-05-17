package file

import (
	"github.com/kazoup/platform/lib/globals"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"time"
)

const (
	DEFAULT_IMAGE_PREVIEW_URL string = globals.SERVER_ADDRESS + "/media/image/http?source=http://www.scaleautomag.com/sitefiles/images/no-preview-available.png"
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
	GetURL() string
	GetExtension() string
	GetModifiedTime() time.Time
	GetContent() string
	GetOptsTimestamps() *OptsKazoupFile
	SetOptsTimestamps(optsKazoupFile *OptsKazoupFile)
	SetHighlight(highlight string)
	SetContentCategory(kazoupCategorization *KazoupCategorization)
	SetEntities(entities *rossetelib.RosseteEntities)
	SetSentiment(sentiment *rossetelib.RosseteSentiment)
}
