package file

import (
	"github.com/kazoup/platform/lib/rossete"
	//googledrive "google.golang.org/api/drive/v3"
	"strings"
	"time"
)

type KazoupGoogleFile struct {
	KazoupFile
	//Original *googledrive.File `json:"original,omitempty"`
}

func (kf *KazoupGoogleFile) PreviewURL(width, height, mode, quality string) string {
	return kf.PreviewUrl
}

func (kf *KazoupGoogleFile) GetID() string {
	return kf.ID
}

func (kf *KazoupGoogleFile) GetName() string {
	return kf.Name
}

func (kf *KazoupGoogleFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupGoogleFile) GetIDFromOriginal() string {
	return kf.OriginalID
}

func (kf *KazoupGoogleFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupGoogleFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupGoogleFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupGoogleFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupGoogleFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupGoogleFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupGoogleFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupGoogleFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupGoogleFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupGoogleFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupGoogleFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupGoogleFile) SetEntities(entities *rossete.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupGoogleFile) SetSentiment(sentiment *rossete.RosseteSentiment) {
	kf.Sentiment = sentiment
}
