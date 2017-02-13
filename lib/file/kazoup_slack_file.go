package file

import (
	"github.com/kazoup/platform/lib/slack"
	"strings"
)

type KazoupSlackFile struct {
	KazoupFile
	Original *slack.SlackFile `json:"original,omitempty"`
}

func (kf *KazoupSlackFile) PreviewURL(width, height, mode, quality string) string {
	// Not all files falling into pictures has all thumbX attrs
	if len(kf.Original.Thumb480) > 0 {
		return kf.Original.Thumb480
	}

	if len(kf.Original.Thumb720) > 0 {
		return kf.Original.Thumb720
	}

	if len(kf.Original.Thumb960) > 0 {
		return kf.Original.Thumb960
	}

	if len(kf.Original.Thumb1024) > 0 {
		return kf.Original.Thumb1024
	}

	return kf.Original.Thumb360
}

func (kf *KazoupSlackFile) GetID() string {
	return kf.ID
}

func (kf *KazoupSlackFile) GetUserID() string {
	return kf.UserId
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

func (kf *KazoupSlackFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupSlackFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupSlackFile) GetBase64() string {
	return ""
}

func (kf *KazoupSlackFile) SetHighlight(s string) {
	kf.Highlight = s
}
