package file

import (
	"fmt"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/globals"
	"strings"
)

type KazoupBoxFile struct {
	KazoupFile
	Original box.BoxFileMeta `json:"original"`
}

func (kf *KazoupBoxFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s%s%s", globals.BoxFileMetadataEndpoint, kf.Original.ID, "/thumbnail.png?min_height=256&min_width=256")

	return url
}

func (kf *KazoupBoxFile) GetID() string {
	return kf.ID
}

func (kf *KazoupBoxFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupBoxFile) GetIDFromOriginal() string {
	return kf.Original.ID
}

func (kf *KazoupBoxFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupBoxFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupBoxFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupBoxFile) GetPathDisplay() string {
	return ""
}

func (kf *KazoupBoxFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupBoxFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupBoxFile) GetBase64() string {
	return ""
}