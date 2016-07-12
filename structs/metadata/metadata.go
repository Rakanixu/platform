package metadata

import (
	"time"

	"github.com/kazoup/platform/structs/intmap"
	"github.com/kazoup/platform/structs/smbattributes"
)

// Metadata model for file
type Metadata struct {
	Mimetype      string                      `json:"mimetype"`
	DocType       string                      `json:"doc_type"`
	DirpathSplit  intmap.Intmap               `json:"dirpath_split"`
	UID           int                         `json:"uid"`
	Extension     string                      `json:"extension"`
	Created       time.Time                   `json:"created"`
	SmbAttributes smbattributes.SmbAttributes `json:"smb_attributes"`
	Modified      time.Time                   `json:"modified"`
	Filename      string                      `json:"filename"`
	Gid           int                         `json:"gid"`
	Dirpath       string                      `json:"dirpath"`
	Accessed      time.Time                   `json:"accessed"`
	Fullpath      string                      `json:"fullpath"`
	Sharepath     string                      `json:"sharepath"`
	FilenameB64   string                      `json:"filename_b64"`
	Size          int64                       `json:"size"`
}
