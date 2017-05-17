package file

import (
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupSlackFile struct {
	KazoupFile
	User     string   `json:"user"`
	Channels []string `json:"channels"`
}

func (kf *KazoupSlackFile) PreviewURL(width, height, mode, quality string) string {
	return kf.PreviewUrl
}

func (kf *KazoupSlackFile) GetID() string {
	return kf.ID
}

func (kf *KazoupSlackFile) GetName() string {
	return kf.Name
}

func (kf *KazoupSlackFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupSlackFile) GetIDFromOriginal() string {
	return kf.OriginalID
}

func (kf *KazoupSlackFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupSlackFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupSlackFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupSlackFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupSlackFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupSlackFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupSlackFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupSlackFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupSlackFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupSlackFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupSlackFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupSlackFile) SetEntities(entities *rossetelib.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupSlackFile) SetSentiment(sentiment *rossetelib.RosseteSentiment) {
	kf.Sentiment = sentiment
}
