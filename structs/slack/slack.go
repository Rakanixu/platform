package slack

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/kazoup/platform/structs"
	"github.com/kazoup/platform/structs/categories"
	"time"
)

// FilesListResponse represents https://slack.com/api/files.list response
type FilesListResponse struct {
	Ok     bool `json:"ok"`
	Files  []SlackFile
	Paging Page `json:"paging"`
}

type Page struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

// SlackFile
type SlackFile struct {
	ID                 string        `json:"id"`
	Created            int           `json:"created"`
	Timestamp          int64         `json:"timestamp"`
	Name               string        `json:"name"`
	Title              string        `json:"title"`
	Mimetype           string        `json:"mimetype"`
	Filetype           string        `json:"filetype"`
	PrettyType         string        `json:"pretty_type"`
	User               string        `json:"user"`
	Editable           bool          `json:"editable"`
	Size               int64         `json:"size"`
	Mode               string        `json:"mode"`
	IsExternal         bool          `json:"is_external"`
	ExternalType       string        `json:"external_type"`
	IsPublic           bool          `json:"is_public"`
	PublicURLShared    bool          `json:"public_url_shared"`
	DisplayAsBot       bool          `json:"display_as_bot"`
	Username           string        `json:"username"`
	URLPrivate         string        `json:"url_private"`
	URLPrivateDownload string        `json:"url_private_download"`
	Thumb64            string        `json:"thumb_64"`
	Thumb80            string        `json:"thumb_80"`
	Thumb360           string        `json:"thumb_360"`
	Thumb360W          int           `json:"thumb_360_w"`
	Thumb360H          int           `json:"thumb_360_h"`
	Thumb480           string        `json:"thumb_480"`
	Thumb480W          int           `json:"thumb_480_w"`
	Thumb480H          int           `json:"thumb_480_h"`
	Thumb160           string        `json:"thumb_160"`
	Thumb720           string        `json:"thumb_720"`
	Thumb720W          int           `json:"thumb_720_w"`
	Thumb720H          int           `json:"thumb_720_h"`
	Thumb960           string        `json:"thumb_960"`
	Thumb960W          int           `json:"thumb_960_w"`
	Thumb960H          int           `json:"thumb_960_h"`
	Thumb1024          string        `json:"thumb_1024"`
	Thumb1024W         int           `json:"thumb_1024_w"`
	Thumb1024H         int           `json:"thumb_1024_h"`
	ImageExifRotation  int           `json:"image_exif_rotation"`
	OriginalW          int           `json:"original_w"`
	OriginalH          int           `json:"original_h"`
	Permalink          string        `json:"permalink"`
	PermalinkPublic    string        `json:"permalink_public"`
	Channels           []string      `json:"channels"`
	Groups             []interface{} `json:"groups"`
	Ims                []interface{} `json:"ims"`
	CommentsCount      int           `json:"comments_count"`
	InitialComment     struct {
		ID        string `json:"id"`
		Created   int    `json:"created"`
		Timestamp int    `json:"timestamp"`
		User      string `json:"user"`
		IsIntro   bool   `json:"is_intro"`
		Comment   string `json:"comment"`
		Channel   string `json:"channel"`
	} `json:"initial_comment"`
}

// NewKazoupFileFromSlackFile constructor
func NewKazoupFileFromSlackFile(s *SlackFile) *structs.KazoupFile {
	t := time.Unix(s.Timestamp, 0)

	return &structs.KazoupFile{
		ID:       getMD5Hash(s.URLPrivate),
		Name:     s.Name,
		URL:      s.URLPrivate,
		Modified: t,
		Size:     s.Size,
		IsDir:    false,
		Category: categories.GetDocType("." + s.Filetype),
		Depth:    0,
		Original: s,
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
