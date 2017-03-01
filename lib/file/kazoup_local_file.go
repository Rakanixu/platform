package file

import (
	"fmt"
	"github.com/kazoup/platform/lib/globals"
	rossetelib "github.com/kazoup/platform/lib/rossete"
	"strings"
	"time"
)

type KazoupLocalFile struct {
	KazoupFile
}

func (kf *KazoupLocalFile) PreviewURL(width, height, mode, quality string) string {
	var url string

	if kf.Category == "Movies" {
		url = fmt.Sprintf("%s/media/frame/%s", globals.SERVER_ADDRESS, strings.TrimLeft(kf.URL, "/local"))
	} else {
		url = fmt.Sprintf("%s/image/local?source=/%s&width=%s&height=%s&mode=%s&quality=%s",
			BASE_URL_FILE_PREVIEW,
			strings.TrimLeft(kf.URL, "/local"),
			width,
			height,
			mode,
			quality,
		)
	}

	return url
}

func (kf *KazoupLocalFile) GetID() string {
	return kf.ID
}

func (kf *KazoupLocalFile) GetUserID() string {
	return kf.UserId
}

func (kf *KazoupLocalFile) GetIDFromOriginal() string {
	return ""
}

func (kf *KazoupLocalFile) GetIndex() string {
	return kf.Index
}

func (kf *KazoupLocalFile) GetDatasourceID() string {
	return kf.DatasourceId
}

func (kf *KazoupLocalFile) GetFileType() string {
	return kf.FileType
}

func (kf *KazoupLocalFile) GetPathDisplay() string {
	return ""
}

func (kf *KazoupLocalFile) GetURL() string {
	return kf.URL
}

func (kf *KazoupLocalFile) GetExtension() string {
	ext := strings.Split(strings.Replace(kf.Name, " ", "-", 1), ".")

	return ext[len(ext)-1]
}

func (kf *KazoupLocalFile) GetBase64() string {
	return ""
}

func (kf *KazoupLocalFile) GetModifiedTime() time.Time {
	return kf.Modified
}

func (kf *KazoupLocalFile) GetContent() string {
	return kf.Content
}

func (kf *KazoupLocalFile) GetOptsTimestamps() *OptsKazoupFile {
	return kf.OptsKazoupFile
}

func (kf *KazoupLocalFile) SetOptsTimestamps(optsKazoupFile *OptsKazoupFile) {
	kf.OptsKazoupFile = optsKazoupFile
}

func (kf *KazoupLocalFile) SetHighlight(s string) {
	kf.Highlight = s
}

func (kf *KazoupLocalFile) SetContentCategory(contentCategory string) {
	kf.ContentCategory = contentCategory
}

func (kf *KazoupLocalFile) SetEntities(entities *rossetelib.RosseteEntities) {
	kf.Entities = entities
}
