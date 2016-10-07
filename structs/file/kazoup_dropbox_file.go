package file

import (
	"github.com/kazoup/platform/structs/dropbox"
	"github.com/kazoup/platform/structs/globals"
)

type KazoupDropboxFile struct {
	KazoupFile
	Original dropbox.DropboxFile
}

func (kf *KazoupDropboxFile) PreviewURL(width, height, mode, quality string) string {
	return ""
}

func (kf *KazoupDropboxFile) GetID() string {
	//return globals.GetMD5Hash(kf.Original.WebViewLink)
	return globals.GetMD5Hash(kf.Name)
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
