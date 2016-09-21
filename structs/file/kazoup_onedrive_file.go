package file

import (
	"github.com/kazoup/platform/structs/globals"
	"github.com/kazoup/platform/structs/onedrive"
)

type KazoupOneDriveFile struct {
	KazoupFile
	Original onedrive.OneDriveFile `json:"original"`
}

func (kf *KazoupOneDriveFile) PreviewURL() string {
	//url := fmt.Sprintf("%s/image/http?source=https://docs.google.com/uc?id=%s", BASE_URL_FILE_PREVIEW, kf.Original.ContentDownloadURL)
	return DEFAULT_IMAGE_PREVIEW_URL
}

func (kf *KazoupOneDriveFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebURL)
}
