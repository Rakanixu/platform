package file

import "time"

// KazoupFile represents all different types
type KazoupFile struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Modified time.Time `json:"modified"`
	FileSize int64     `json:"file_size"`
	IsDir    bool      `json:"is_dir"`
	Category string    `json:"category"`
	Depth    int64     `json:"depth"`
	FileType string    `json:"file_type"`
}

func (kf *KazoupFile) PreviewURL() string {
	return DEFAULT_IMAGE_PREVIEW_URL
}
