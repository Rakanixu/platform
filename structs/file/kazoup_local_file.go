package file

import (
	"fmt"
	"strings"
)

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL() string {
	url := fmt.Sprintf("%s/image/local?source=/%s", BASE_URL_FILE_PREVIEW, strings.TrimLeft(kf.URL, "/local"))
	return url
}

func (kf *KazoupLocalFile) GetID() string {
	return kf.ID
}
