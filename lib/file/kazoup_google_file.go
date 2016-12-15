package file

import (
	"github.com/kazoup/platform/lib/globals"
	googledrive "google.golang.org/api/drive/v3"
	"strings"
)

type KazoupGoogleFile struct {
	KazoupFile
	Original *googledrive.File `json:"original,omitempty"`
}

func (kf *KazoupGoogleFile) PreviewURL(width, height, mode, quality string) string {
	return ""
}

func (kf *KazoupGoogleFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebViewLink)
}

func (kf *KazoupGoogleFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupGoogleFile) GetIDFromOriginal() string {
	return kf.Original.Id
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

func (kf *KazoupGoogleFile) GetPathDisplay() string {
	return ""
}

func (kf *KazoupGoogleFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupGoogleFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupGoogleFile) GetBase64() string {
	return ""
}
