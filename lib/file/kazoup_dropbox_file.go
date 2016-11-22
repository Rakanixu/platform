package file

import (
	"github.com/kazoup/platform/lib/dropbox"
	"strings"
)

type KazoupDropboxFile struct {
	KazoupFile
	Original dropbox.DropboxFile `json:"original"`
}

func (kf *KazoupDropboxFile) PreviewURL(width, height, mode, quality string) string {
	return ""
}

func (kf *KazoupDropboxFile) GetID() string {
	return kf.ID
}

func (kf *KazoupDropboxFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupDropboxFile) GetIDFromOriginal() string {
	return kf.Original.ID
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

func (kf *KazoupDropboxFile) GetPathDisplay() string {
	return kf.Original.PathDisplay
}

func (kf *KazoupDropboxFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupDropboxFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupDropboxFile) GetBase64() string {
	return ""
}