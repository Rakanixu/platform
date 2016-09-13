package file

import (
	"fmt"
	googledrive "google.golang.org/api/drive/v3"
)

type KazoupGoogleFile struct {
	KazoupFile
	Original googledrive.File `json:"original"`
}

func (kf *KazoupGoogleFile) PreviewURL() string {
	url := fmt.Sprintf("%s/image/http?source=%s", BASE_URL_FILE_PREVIEW, kf.Original.ThumbnailLink)
	return url
}
