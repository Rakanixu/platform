package file

import "time"

// KazoupFile represents all different types
type KazoupFile struct {
	ID           string    `json:"id"`
	UserId       string    `json:"user_id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Modified     time.Time `json:"modified"`
	FileSize     int64     `json:"file_size"`
	IsDir        bool      `json:"is_dir"`
	Category     string    `json:"category"`
	MimeType     string    `json:"mime_type"`
	Depth        int64     `json:"depth"`
	FileType     string    `json:"file_type"`
	LastSeen     int64     `json:"last_seen"`
	Access       string    `json:"access"`
	DatasourceId string    `json:"datasource_id"`
	Index        string    `json:"index,omitempty"`   // Index the file will be pushed to
	Content      string    `json:"content,omitempty"` // Content extarcted from tika
}
