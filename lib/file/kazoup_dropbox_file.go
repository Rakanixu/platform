package file

import (
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupDropboxFile struct {
	KazoupFile
	DropboxUsers    []dropbox.DropboxUser    `json:"dropbox_users,omitempty"`
	DropboxInvitees []dropbox.DropboxInvitee `json:"dropbox_invitees,omitempty"`
	//Original *dropbox.DropboxFile `json:"original,omitempty"`
}

func (kf *KazoupDropboxFile) PreviewURL(width, height, mode, quality string) string {
	return kf.PreviewUrl
}

func (kf *KazoupDropboxFile) GetID() string {
	return kf.ID
}

func (kf *KazoupDropboxFile) GetName() string {
	return kf.Name
}

func (kf *KazoupDropboxFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupDropboxFile) GetIDFromOriginal() string {
	return kf.OriginalID
}

func (kf *KazoupDropboxFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupDropboxFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupDropboxFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupDropboxFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupDropboxFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupDropboxFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupDropboxFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupDropboxFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupDropboxFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupDropboxFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupDropboxFile) SetContentCategory(c *KazoupCategorization) {
	kf.KazoupCategorization = c
}

func (kf *KazoupDropboxFile) SetEntities(entities *rossete.RosseteEntities) {
	kf.Entities = entities
}

func (kf *KazoupDropboxFile) SetSentiment(sentiment *rossete.RosseteSentiment) {
	kf.Sentiment = sentiment
}
