package file

import "time"

// KazoupFile represents all different types
type KazoupFile struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Modified time.Time `json:"modified"`
	Size     int64     `json:"size"`
	IsDir    bool      `json:"is_dir"`
	Category string    `json:"category"`
	Depth    int64     `json:"depth"`
	FileType string    `json:"file_type"`
}

func (kf *KazoupFile) PreviewURL() string {
	return kf.URL
}
