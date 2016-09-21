package file

import (
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	googledrive "google.golang.org/api/drive/v3"
)

type KazoupGoogleFile struct {
	KazoupFile
	Original googledrive.File `json:"original"`
}

func (kf *KazoupGoogleFile) PreviewURL() string {
	size := "s600"
	link := fmt.Sprintf("%s%s", kf.Original.ThumbnailLink[:len(kf.Original.ThumbnailLink)-4], size)
	url := fmt.Sprintf("%s/image/http?source=%s", BASE_URL_FILE_PREVIEW, link)
	return url
}

func (kf *KazoupGoogleFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebViewLink)
}
