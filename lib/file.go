package structs

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/intmap"
)

const (
	DEFAULT_IMAGE_URL string = "http://placehold.it/350x150"
)

// OriginalFile interface
type OriginalFile interface {
}

// KazoupFile represents all different types
type KazoupFile struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	URL      string       `json:"url"`
	Modified time.Time    `json:"modified"`
	Size     int64        `json:"size"`
	IsDir    bool         `json:"is_dir"`
	Category string       `json:"category"`
	Depth    int64        `json:"depth"`
	FileType string       `json:"file_type"`
	Original OriginalFile `json:"original,omitempty"`
}

func NewKazoupFileFromString(s string) (*KazoupFile, error) {
	f := &KazoupFile{}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, f); err != nil {
		return nil, err
	}
	return f, nil
}

// DesktopFile ...
type DesktopFile struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	URL      string      `json:"url"`
	Modified time.Time   `json:"modified"`
	Size     int64       `json:"size"`
	IsDir    bool        `json:"is_dir"`
	Mode     os.FileMode `json:"mode"`
	Category string      `json:"category"`
	Depth    int64       `json:"depth"`
}

// Desktop file optimized
type DesktopFileOptimised struct {
	N string
	U string
	M time.Time
	S int64
	D bool
	P os.FileMode
}

// LocalFile model
type LocalFile struct {
	Type string
	Path string
	Info os.FileInfo
	Id   string
}

func NewDesktopFile(lf *LocalFile) *DesktopFile {
	return &DesktopFile{
		ID:       lf.Id,
		Name:     lf.Info.Name(),
		URL:      "/local" + lf.Path,
		Modified: lf.Info.ModTime(),
		Size:     lf.Info.Size(),
		IsDir:    lf.Info.IsDir(),
		Mode:     lf.Info.Mode(),
		Category: categories.GetDocType(filepath.Ext(lf.Info.Name())),
		Depth:    UrlDepth(lf.Path),
	}
}

// MockFile model
type MockFile struct {
	Filename     string
	Mimetype     string
	Dirpath      string
	Fullpath     string
	DirpathSplit intmap.Intmap
	Sharepath    string
	Extension    string
	DocType      string
}

var directories = [...]string{"aaa", "bbb", "ccc", "ddd", "eee"}
var extensions = [...]string{".js", ".go", ".png", ".avi", ".txt"}
var docTypes = [...]string{"JavaScript", "Golang", "Images", "Videos", "Documents"}
var mimeTypes = [...]string{"application/javascript", "application", "image/png", "video/avi", "text/plain"}

// GenerateData for mock file
func (mf *MockFile) GenerateData() {
	index := randomdata.Number(0, 4)
	path := "//127.0.0.1/"

	mf.Filename += extensions[index]
	mf.Extension = extensions[index]
	mf.DocType = docTypes[index]
	mf.Mimetype = mimeTypes[index]

	for i := 0; i < index; i++ {
		path += directories[randomdata.Number(0, 4)] + "/"
	}
	mf.Dirpath = path
	path += mf.Filename

	mf.Fullpath = path
	mf.DirpathSplit = pathToIntmap(path)
	mf.Sharepath = "//vol1/"
}

// NewMockFile constructor
func NewMockFile() *DesktopFile {
	return &DesktopFile{}
}

func pathToIntmap(path string) intmap.Intmap {
	results := make(intmap.Intmap)
	dir := filepath.Dir(path)
	parts := strings.Split(dir, "/")
	for k, v := range parts {
		if k == 0 {
			continue
		} else if k == 1 {
			results[k-1] = "//" + v
		} else {
			results[k-1] = "/" + filepath.Join(results[k-2], v)
		}
	}
	return results
}

func UrlDepth(str string) int64 {
	return int64(len(strings.Split(str, "/")) - 1)
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
