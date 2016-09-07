package googledrive

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/kazoup/platform/structs"
	"github.com/kazoup/platform/structs/categories"
	googledrive "google.golang.org/api/drive/v3"
	"time"
)

// NewKazoupFileFromSlackFile constructor
func NewKazoupFileFromGoogleDriveFile(g *googledrive.File) *structs.KazoupFile {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", g.ModifiedTime)

	return &structs.KazoupFile{
		ID:       getMD5Hash(g.WebViewLink),
		Name:     g.Name,
		URL:      g.WebViewLink,
		Modified: t,
		Size:     g.Size,
		IsDir:    false,
		Category: categories.GetDocType("." + g.FullFileExtension),
		Depth:    0,
		Original: g,
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
