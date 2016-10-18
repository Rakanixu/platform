package file

import (
	"github.com/kazoup/platform/structs/dropbox"
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
