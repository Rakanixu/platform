package file

import (
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"time"
)

type OptsKazoupFile struct {
	ContentTimestamp           *time.Time `json:"content_timestamp,omitempty"`            // Last time content was extracted
	TagsTimestamp              *time.Time `json:"tags_timestamp,omitempty"`               // Last time tags was extracted
	AudioTimestamp             *time.Time `json:"audio_timestamp,omitempty"`              // Last time audio was extracted
	TextAnalyzedTimestamp      *time.Time `json:"text_analyzed_timestamp,omitempty"`      // Last time text analytics was runned over content (Rossete)
	SentimentAnalyzedTimestamp *time.Time `json:"sentiment_analyzed_timestamp,omitempty"` // Last time sentiment analytics was runned over content (Rossete)
	ThumbnailTimestamp         *time.Time `json:"thumbnail_timestamp,omitempty"`          // Last time thumbnail was generated
}

type KazoupCategorization struct {
	ContentCategory *string `json:"content_category"`
}

// KazoupFile represents all different types
type KazoupFile struct {
	ID                   string                       `json:"id"`
	OriginalID           string                       `json:"original_id"`
	OriginalDownloadRef  string                       `json:"original_download_ref"`
	PreviewUrl           string                       `json:"preview_url"`
	UserId               string                       `json:"user_id"`
	Name                 string                       `json:"name"`
	URL                  string                       `json:"url"`
	Modified             time.Time                    `json:"modified"`
	FileSize             int64                        `json:"file_size"`
	IsDir                bool                         `json:"is_dir"`
	Category             string                       `json:"category"`
	MimeType             string                       `json:"mime_type"`
	Depth                int64                        `json:"depth"`
	FileType             string                       `json:"file_type"`
	LastSeen             int64                        `json:"last_seen"`
	Access               string                       `json:"access"`
	DatasourceId         string                       `json:"datasource_id"`
	Index                string                       `json:"index,omitempty"`     // Index the file will be pushed to
	Content              string                       `json:"content,omitempty"`   // Content extarcted from tika
	Highlight            string                       `json:"highlight,omitempty"` // Highlight search term
	Tags                 []string                     `json:"tags,omitempty"`      // Tags from cloud vision or other tags releted with file content
	OptsKazoupFile       *OptsKazoupFile              `json:"opts_kazoup_file,omitempty"`
	KazoupCategorization *KazoupCategorization        `json:"kazoup_categorization,omitempty"`
	Entities             *rossetelib.RosseteEntities  `json:"entities,omitempty"`
	Sentiment            *rossetelib.RosseteSentiment `json:"sentiment,omitempty"`
}
