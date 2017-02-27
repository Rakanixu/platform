package file

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	gmail "github.com/kazoup/platform/lib/gmail"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"strings"
)

type KazoupGmailFile struct {
	KazoupFile
	Original *gmail.GmailFile `json:"original,omitempty"`
}

func (kf *KazoupGmailFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s/%s", globals.GmailEndpoint, kf.Original.Id)

	return url
}

func (kf *KazoupGmailFile) GetID() string {
	return kf.ID
}

func (kf *KazoupGmailFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupGmailFile) GetIDFromOriginal() string {
	return kf.Original.Id
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

func (kf *KazoupGmailFile) GetPathDisplay() string {
	return ""
}

func (kf *KazoupGmailFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupGmailFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupGmailFile) GetBase64() string {
	return kf.Original.Base64
}

func (kf *KazoupGmailFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupGmailFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupGmailFile) SetEntities(entities *rossetelib.RosseteEntities) {
	kf.Entities = entities
}
