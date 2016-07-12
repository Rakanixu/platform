package content

import (
	"github.com/kazoup/platform/structs/countvalue"
)

// Content model for file
type Content struct {
	ContentTime   float64               `json:"content_time"`
	TikaMetadata  string                `json:"tika_metadata"`
	Unique        bool                  `json:"unique"`
	ContentStored bool                  `json:"content_stored"`
	Text          string                `json:"text"`
	Checksum      string                `json:"checksum"`
	ChecksumTime  float64               `json:"checksum_time"`
	Money         countvalue.CountValue `json:"money"`
	Organisation  countvalue.CountValue `json:"organisation"`
	Percent       countvalue.CountValue `json:"percent"`
	Person        countvalue.CountValue `json:"person"`
	Location      countvalue.CountValue `json:"location"`
	Time          countvalue.CountValue `json:"time"`
	Date          countvalue.CountValue `json:"date"`
}
