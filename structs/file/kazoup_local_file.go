package file

import (
	"fmt"
	"github.com/kazoup/platform/structs/globals"
	"strings"
)

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL(width, height, mode, quality string) string {
	var url string

	if kf.Category == "Movies" {
		url = fmt.Sprintf("%s/media/frame/%s", globals.SERVER_ADDRESS, strings.TrimLeft(kf.URL, "/local"))
	} else {
		url = fmt.Sprintf("%s/image/local?source=/%s&width=%s&height=%s&mode=%s&quality=%s",
			BASE_URL_FILE_PREVIEW,
			strings.TrimLeft(kf.URL, "/local"),
			width,
			height,
			mode,
			quality,
		)
	}

	return url
}

func (kf *KazoupLocalFile) GetID() string {
	return kf.ID
}

func (kf *KazoupLocalFile) GetIDFromOriginal() string {
	return ""
}

func (kf *KazoupLocalFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupLocalFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupLocalFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupLocalFile) GetPathDisplay() string {
	return ""
}
