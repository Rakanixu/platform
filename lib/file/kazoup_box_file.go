package file

import (
	"fmt"
	//"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupBoxFile struct {
	KazoupFile
	//Original *box.BoxFileMeta `json:"original,omitempty"`
}

func (kf *KazoupBoxFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s%s%s", globals.BoxFileMetadataEndpoint, kf.OriginalID, "/thumbnail.png?min_height=256&min_width=256")

	return url
}

func (kf *KazoupBoxFile) GetID() string {
	return kf.ID
}

func (kf *KazoupBoxFile) GetName() string {
	return kf.Name
}

func (kf *KazoupBoxFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupBoxFile) GetIDFromOriginal() string {
	return kf.OriginalID
}

func (kf *KazoupBoxFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupBoxFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupBoxFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupBoxFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupBoxFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupBoxFile) GetBase64() string {
	return ""
}

func (kf *KazoupBoxFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupBoxFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupBoxFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupBoxFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupBoxFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupBoxFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupBoxFile) SetEntities(entities *rossete.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupBoxFile) SetSentiment(sentiment *rossete.RosseteSentiment) {
	kf.Sentiment = sentiment
}
