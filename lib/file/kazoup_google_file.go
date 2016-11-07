package file

import (
	"github.com/kazoup/platform/lib/globals"
	googledrive "google.golang.org/api/drive/v3"
)

type KazoupGoogleFile struct {
	KazoupFile
	Original googledrive.File `json:"original"`
}

func (kf *KazoupGoogleFile) PreviewURL(width, height, mode, quality string) string {
	return ""
}

func (kf *KazoupGoogleFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebViewLink)
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
