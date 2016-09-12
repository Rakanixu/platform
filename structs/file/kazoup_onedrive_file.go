package file

import "github.com/kazoup/platform/structs/onedrive"

type KazoupOneDriveFile struct {
	KazoupFile
	Original onedrive.OneDriveFile `json:"original"`
}

func (kf *KazoupOneDriveFile) PreviewURL() string {
	return DEFAULT_IMAGE_PREVIEW_URL
}
