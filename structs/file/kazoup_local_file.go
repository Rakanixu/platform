package file

import (
	"fmt"
	"strings"
)

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s/image/local?source=/%s&width=%s&height=%s&mode=%s&quality=%s",
		BASE_URL_FILE_PREVIEW,
		strings.TrimLeft(kf.URL, "/local"),
		width,
		height,
		mode,
		quality,
	)

	return url
}

func (kf *KazoupLocalFile) GetID() string {
	return kf.ID
}

func (kf *KazoupLocalFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupLocalFile) GetDatasourceID() string {
	return kf.DatasourceId
}
