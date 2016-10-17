package file

import (
	"github.com/kazoup/platform/structs/box"
	"github.com/kazoup/platform/structs/globals"
	"fmt"
)

type KazoupBoxFile struct {
	KazoupFile
	Original box.BoxFileMeta
}

func (kf *KazoupBoxFile) PreviewURL(width, height, mode, quality string) string {
	url := fmt.Sprintf("%s%s%s", globals.BoxFileMetadataEndpoint, kf.Original.ID, "/thumbnail.png?min_height=256&min_width=256")

	return url
}

func (kf *KazoupBoxFile) GetID() string {
	return kf.ID
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
