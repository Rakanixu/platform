package file

import "github.com/kazoup/platform/structs/slack"

type KazoupSlackFile struct {
	KazoupFile
	Original slack.SlackFile `json:"original"`
}

func (kf *KazoupSlackFile) PreviewURL() string {
	return kf.Original.Permalink
}
