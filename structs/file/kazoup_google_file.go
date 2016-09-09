package file

import (
//googledrive "google.golang.org/api/drive/v3"
)

type KazoupGoogleFile struct {
	KazoupFile
	//Original googledrive.File `json:"original"`
}

func (kf *KazoupGoogleFile) PreviewURL() string {
	return kf.URL
}
