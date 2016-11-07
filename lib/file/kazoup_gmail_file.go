package file

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	gmail "google.golang.org/api/gmail/v1"
)

type KazoupGmailFile struct {
	KazoupFile
	Original gmail.Message `json:"original"`
}

func (kf *KazoupGmailFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s/%s", globals.GmailEndpoint, kf.Original.Id)

	return url
}

func (kf *KazoupGmailFile) GetID() string {
	return kf.ID
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
