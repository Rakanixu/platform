package file

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/kazoup/platform/structs/categories"
	local "github.com/kazoup/platform/structs/local"
	slack "github.com/kazoup/platform/structs/slack"
	googledrive "google.golang.org/api/drive/v3"
)

type File interface {
	GetByID(string) (File, error)
	PreviewURL() string
}

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

type KazoupSlackFile struct {
	KazoupFile
	slack.SlackFile
}

func (kf *KazoupSlackFile) PreviewURL() string {
	return kf.Permalink
}

type KazoupGoogleFile struct {
	KazoupFile
	googledrive.File
}

func (kf *KazoupGoogleFile) PreviewURL() string {
	return kf.ThumbnailLink
}

type KazoupLocalFile struct {
	KazoupFile
	local.LocalFile
}

func (kf *KazoupLocalFile) PreviewURL() string {
	return kf.URL
}

func NewFileFromString(s string) (File, error) {
	kf := &KazoupFile{}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, kf); err != nil {
		return nil, err
	}
	switch kf.FileType {
	case "local":
		f, err := NewLocalFileFromString(string)
		if err != nil {
			return nil, err
		}
		return f, nil
	default:
		return nil, errors.New("Hmmm crap that should not happen")
	}
}

func NewLocalFileFromString(string) (*LocalFile, err) {
}

func NewKazoupSlackFile(sf *slack.SlackFile) *KazoupSlackFile {
	f := &KazoupFile{
		ID:       sf.ID,
		Name:     sf.Name,
		URL:      sf.Permalink,
		Modified: time.Unix(sf.Timestamp, 0),
		Size:     sf.Size,
		IsDir:    false,
		Category: categories.GetDocType("." + sf.Filetype),
	}
	return &KazoupSlackFile{*f, *sf}
}
