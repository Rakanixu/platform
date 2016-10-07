package file

import "github.com/kazoup/platform/structs/slack"

type KazoupSlackFile struct {
	KazoupFile
	Original slack.SlackFile `json:"original"`
}

func (kf *KazoupSlackFile) PreviewURL(width, height, mode, quality string) string {
	return kf.Original.Thumb480
}

func (kf *KazoupSlackFile) GetID() string {
	return kf.ID
}

func (kf *KazoupSlackFile) GetIDFromOriginal() string {
	return kf.Original.ID
}

func (kf *KazoupSlackFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupSlackFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupSlackFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupSlackFile) GetPathDisplay() string {
	return ""
}
