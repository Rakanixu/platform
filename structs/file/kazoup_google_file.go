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

func (kf *KazoupGoogleFile) PreviewURL(width, height, mode, quality string) string {
	sz := "w240"

	// Prioritize width over height if both defined
	if len(height) > 0 {
		sz = "h" + height
	}
	if len(width) > 0 {
		sz = "w" + width
	}

	url := fmt.Sprintf("%s&sz=%s&id=%s", GOOGLE_DRIVE_THUMBNAIL, sz, kf.Original.Id)

	//https: //drive.google.com/thumbnail?authuser=0&sz=w320&id=ee8366992d921d448d297e926638ecea

	return url
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
