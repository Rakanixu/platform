package structs

import (
	"crypto/rand"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/kazoup/platform/structs/content"
	"github.com/kazoup/platform/structs/intmap"
	"github.com/kazoup/platform/structs/metadata"
	"github.com/kazoup/platform/structs/permissions"
)

// File model
type File struct {
	ExistsOnDisk bool `json:"exists_on_disk"`
	//ID              string                  `json:"_id"`
	ArchiveComplete bool                    `json:"archive_complete"`
	FirstSeen       time.Time               `json:"first_seen"`
	IDB64           string                  `json:"id_b64"`
	LastSeen        time.Time               `json:"last_seen"`
	Content         content.Content         `json:"content"`
	Metadata        metadata.Metadata       `json:"metadata"`
	Permissions     permissions.Permissions `json:"permissions"`
}

// LocalFile model
type LocalFile struct {
	Type string
	Path string
	Info os.FileInfo
}

// NewFileFromLocal file constructor
func NewFileFromLocal(lf *LocalFile) *File {
	return &File{
		//ID:              "/" + lf.Path + ":" + strconv.FormatInt(lf.Info.ModTime().Unix(), 10),
		ExistsOnDisk: true,
		//	ID:              pseudoUUID(),
		ArchiveComplete: false,
		FirstSeen:       time.Now(),
		Content:         content.Content{},
		Metadata: metadata.Metadata{
			Mimetype:     mime.TypeByExtension(filepath.Ext(lf.Info.Name())),
			DirpathSplit: pathToIntmap(lf.Path),
			Extension:    filepath.Ext(lf.Info.Name()),
			Created:      lf.Info.ModTime(),
			Modified:     lf.Info.ModTime(),
			Filename:     lf.Info.Name(),
			Dirpath:      filepath.Dir(lf.Path),
			Accessed:     lf.Info.ModTime(),
			Fullpath:     lf.Path,
			Sharepath:    filepath.VolumeName(lf.Path),
			Size:         lf.Info.Size(),
		},
		Permissions: permissions.Permissions{},
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
}

var directories = [...]string{"aaa", "bbb", "ccc", "ddd", "eee"}
var extensions = [...]string{".js", ".go", ".png", ".avi", ".txt"}
var mimeTypes = [...]string{"application/javascript", "application", "image/png", "video/avi", "text/plain"}

// GenerateData for mock file
func (mf *MockFile) GenerateData() {
	index := randomdata.Number(0, 4)
	path := "/127.0.0.1/"

	mf.Filename += extensions[index]
	mf.Extension = extensions[index]
	mf.Mimetype = mimeTypes[index]

	for i := 0; i < index; i++ {
		path += directories[randomdata.Number(0, 4)] + "/"
	}
	mf.Dirpath = path
	path += mf.Filename

	mf.Fullpath = path
	mf.DirpathSplit = pathToIntmap(path)
	mf.Sharepath = "/vol1/"
}

// NewMockFile constructor
func NewMockFile() *File {
	mockFile := &MockFile{
		Filename: randomdata.SillyName(),
	}
	mockFile.GenerateData()

	return &File{
		ExistsOnDisk: true,
		//ID:              pseudoUUID(),
		ArchiveComplete: randomdata.Boolean(),
		FirstSeen: time.Date(
			randomdata.Number(1990, 2015),
			time.November,
			randomdata.Number(1, 28),
			0, 0, 0, 0, time.UTC,
		),
		Content: content.Content{},
		Metadata: metadata.Metadata{
			Mimetype:     mockFile.Mimetype,
			DirpathSplit: mockFile.DirpathSplit,
			Extension:    mockFile.Extension,
			Created: time.Date(
				randomdata.Number(1990, 2015),
				time.November,
				randomdata.Number(1, 28),
				0, 0, 0, 0, time.UTC,
			),
			Modified: time.Date(
				randomdata.Number(1995, 2016),
				time.January,
				randomdata.Number(1, 28),
				0, 0, 0, 0, time.UTC,
			),
			Accessed: time.Date(
				randomdata.Number(2000, 2016),
				time.January,
				randomdata.Number(1, 28),
				0, 0, 0, 0, time.UTC,
			),
			Filename:  mockFile.Filename,
			Dirpath:   mockFile.Dirpath,
			Fullpath:  mockFile.Fullpath,
			Sharepath: mockFile.Sharepath,
			Size:      int64(randomdata.Number(1024, 1048576)),
		},
		Permissions: permissions.Permissions{},
	}
}

func pathToIntmap(path string) intmap.Intmap {
	results := make(intmap.Intmap)
	dir := filepath.Dir(path)
	parts := strings.Split(dir, "/")
	for k, v := range parts {
		if k == 0 {
			results[k] = "/" + v

		} else {
			results[k] = filepath.Join(results[k-1], v)
		}
	}
	return results
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
