package file

import "github.com/kazoup/platform/structs/slack"

type KazoupSlackFile struct {
	KazoupFile
	Original slack.SlackFile `json:"original"`
}

func (kf *KazoupSlackFile) PreviewURL() string {
	return kf.Original.PermalinkPublic
}

func (kf *KazoupSlackFile) GetID() string {
	return kf.ID
}

func (kf *KazoupSlackFile) GetIndex() string {
	return kf.Index
}
