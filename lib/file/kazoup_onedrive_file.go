package file

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/onedrive"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupOneDriveFile struct {
	KazoupFile
	Original *onedrive.OneDriveFile `json:"original,omitempty"`
}

func (kf *KazoupOneDriveFile) PreviewURL(width, height, mode, quality string) string {
	return DEFAULT_IMAGE_PREVIEW_URL
}

func (kf *KazoupOneDriveFile) GetID() string {
	return globals.GetMD5Hash(kf.Original.WebURL)
}

func (kf *KazoupOneDriveFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupOneDriveFile) GetIDFromOriginal() string {
	return kf.Original.ID
}

func (kf *KazoupOneDriveFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupOneDriveFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupOneDriveFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupOneDriveFile) GetPathDisplay() string {
	return ""
}

func (kf *KazoupOneDriveFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupOneDriveFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupOneDriveFile) GetBase64() string {
	return ""
}

func (kf *KazoupOneDriveFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupOneDriveFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupOneDriveFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupOneDriveFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupOneDriveFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupOneDriveFile) SetContentCategory(contentCategory string) {
	kf.ContentCategory = contentCategory
}

func (kf *KazoupOneDriveFile) SetEntities(entities *rossetelib.RosseteEntities) {
	kf.Entities = entities
}
