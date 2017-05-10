package file

import (
	"github.com/kazoup/platform/lib/rossete"
	"time"
)

type KazoupMockFile struct {
	KazoupFile
	Original interface{} `json:"original,omitempty"`
}

func (kf *KazoupMockFile) PreviewURL(width, height, mode, quality string) string {
	return "PreviewURL"
}

func (kf *KazoupMockFile) GetID() string {
	return "GetID"
}

func (kf *KazoupMockFile) GetName() string {
	return "GetName"
}

func (kf *KazoupMockFile) GetUserID() string {
	return "GetUserID"
}

func (kf *KazoupMockFile) GetIDFromOriginal() string {
	return "GetIDFromOriginal"
}

func (kf *KazoupMockFile) GetIndex() string {
	return "GetIndex"
}

func (kf *KazoupMockFile) GetDatasourceID() string {
	return "GetDatasourceID"
}

func (kf *KazoupMockFile) GetFileType() string {
	return "GetFileType"
}

func (kf *KazoupMockFile) GetPathDisplay() string {
	return "GetPathDisplay"
}

func (kf *KazoupMockFile) GetURL() string {
	return "GetURL"
}

func (kf *KazoupMockFile) GetExtension() string {
	return "GetExtension"
}

func (kf *KazoupMockFile) GetBase64() string {
	return "GetBase64"
}

func (kf *KazoupMockFile) GetModifiedTime() time.Time {
	return time.Now()
}

func (kf *KazoupMockFile) GetContent() string {
	return "GetContent"
}

func (kf *KazoupMockFile) GetOptsTimestamps() *OptsKazoupFile {
	return &OptsKazoupFile{}
}

func (kf *KazoupMockFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupMockFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupMockFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupMockFile) SetEntities(entities *rossete.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupMockFile) SetSentiment(sentiment *rossete.RosseteSentiment) {
	kf.Sentiment = sentiment
}
