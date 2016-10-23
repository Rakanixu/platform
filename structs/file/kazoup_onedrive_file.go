package file

import (
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
)

type KazoupOneDriveFile struct {
	KazoupFile
	Original onedrive.OneDriveFile `json:"original"`
}

func (kf *KazoupOneDriveFile) PreviewURL(width, height, mode, quality string) string {
	return DEFAULT_IMAGE_PREVIEW_URL
}

func (kf *KazoupOneDriveFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebURL)
}

func (kf *KazoupOneDriveFile) GetIDFromOriginal() string {
	return kf.Original.ID
}

func (kf *KazoupOneDriveFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupOneDriveFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupOneDriveFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupOneDriveFile) GetPathDisplay() string {
	return ""
}
