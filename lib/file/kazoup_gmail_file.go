package file

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	//"github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupGmailFile struct {
	KazoupFile
	MessageId    string `json:"message_id,omitempty"`
	AttachmentId string `json:"attachment_id,omitempty"`
	//Original *gmail.GmailFile `json:"original,omitempty"`
}

func (kf *KazoupGmailFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s%s", globals.GmailEndpoint, kf.OriginalID)

	return url
}

func (kf *KazoupGmailFile) GetID() string {
	return kf.ID
}

func (kf *KazoupGmailFile) GetName() string {
	return kf.Name
}

func (kf *KazoupGmailFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupGmailFile) GetIDFromOriginal() string {
	return kf.OriginalID
}

func (kf *KazoupGmailFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupGmailFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupGmailFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupGmailFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupGmailFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupGmailFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupGmailFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupGmailFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupGmailFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupGmailFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupGmailFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupGmailFile) SetEntities(entities *rossete.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupGmailFile) SetSentiment(sentiment *rossete.RosseteSentiment) {
	kf.Sentiment = sentiment
}
