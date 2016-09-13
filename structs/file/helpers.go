package file

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/kazoup/platform/structs/categories"
	"github.com/kazoup/platform/structs/local"
	"github.com/kazoup/platform/structs/onedrive"
	"github.com/kazoup/platform/structs/slack"
	"github.com/micro/go-micro/client"

	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/structs/globals"
	googledrive "google.golang.org/api/drive/v3"
)

func GetFileByID(id string) (File, error) {
	dbclient := db.NewDBClient("go.micro.srv.db", client.NewClient())

	// get file URL from DB
	dbreq := db.ReadRequest{
		Index: "files",
		Type:  "file",
		Id:    id,
	}
	dbres, err := dbclient.Read(context.TODO(), &dbreq)
	if err != nil {
		return nil, err
	}
	f, err := NewFileFromString(dbres.Result)
	return f, err

}

func NewFileFromString(s string) (File, error) {
	kf := &KazoupFile{}
	if err := json.Unmarshal([]byte(s), kf); err != nil {
		return nil, errors.New("Error unmarsahling NewFileFromString")
	}

	switch kf.FileType {
	case "local":
		klf := &KazoupLocalFile{}
		if err := json.Unmarshal([]byte(s), klf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case local")
		}
		return klf, nil
	case "slack":
		ksf := &KazoupSlackFile{}
		if err := json.Unmarshal([]byte(s), ksf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case slack")
		}
		return ksf, nil
	case "googledrive":
		kgf := &KazoupGoogleFile{}
		log.Printf("Googledrive string : %s \n", s)
		if err := json.Unmarshal([]byte(s), kgf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case googledrive")
		}
		return kgf, nil
	case "onedrive":
		of := &onedrive.OneDriveFile{}
		if err := json.Unmarshal([]byte(s), of); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case onedrive")
		}
		return &KazoupOneDriveFile{*kf, *of}, nil
	default:
		return nil, errors.New("Error constructing file type")
	}
}

// NewKazoupFileFromGoogleDriveFile constructor
func NewKazoupFileFromGoogleDriveFile(g *googledrive.File) *KazoupGoogleFile {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", g.ModifiedTime)

	kf := &KazoupFile{
		ID:       getMD5Hash(g.WebViewLink),
		Name:     g.Name,
		URL:      g.WebViewLink,
		Modified: t,
		FileSize: g.Size,
		IsDir:    false,
		Category: categories.GetDocType("." + g.FullFileExtension),
		Depth:    0,
		FileType: globals.GoogleDrive,
	}
	return &KazoupGoogleFile{*kf, *g}
}

// NewKazoupFileFromSlackFile constructor
func NewKazoupFileFromSlackFile(s *slack.SlackFile) *KazoupSlackFile {
	t := time.Unix(s.Timestamp, 0)

	kf := &KazoupFile{
		ID:       getMD5Hash(s.URLPrivate),
		Name:     s.Name,
		URL:      s.URLPrivate,
		Modified: t,
		FileSize: s.Size,
		IsDir:    false,
		Category: categories.GetDocType("." + s.Filetype),
		Depth:    0,
		FileType: globals.Slack,
	}
	return &KazoupSlackFile{*kf, *s}
}

func NewKazoupFileFromLocal(lf *local.LocalFile) *KazoupLocalFile {
	// don;t save all LocalFile as mmost of data is same as KazoupFile just pass file mode
	kf := &KazoupFile{
		ID:       getMD5Hash(lf.Path),
		Name:     lf.Info.Name(),
		URL:      "/local" + lf.Path,
		Modified: lf.Info.ModTime(),
		FileSize: lf.Info.Size(),
		IsDir:    lf.Info.IsDir(),
		Category: categories.GetDocType(filepath.Ext(lf.Info.Name())),
		Depth:    UrlDepth(lf.Path),
		FileType: globals.Local,
	}
	return &KazoupLocalFile{*kf}

}

// NewKazoupFileFromOneDriveFile constructor
func NewKazoupFileFromOneDriveFile(o *onedrive.OneDriveFile) *KazoupOneDriveFile {
	isDir := true
	name := strings.Split(o.Name, ".")

	if len(o.File.MimeType) > 0 {
		isDir = false
	}

	kf := &KazoupFile{
		ID:       getMD5Hash(o.WebURL),
		Name:     o.Name,
		URL:      o.WebURL,
		Modified: o.LastModifiedDateTime,
		FileSize: o.Size,
		IsDir:    isDir,
		Category: categories.GetDocType("." + name[len(name)-1]),
		Depth:    0,
		FileType: globals.OneDrive,
	}
	return &KazoupOneDriveFile{*kf, *o}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func UrlDepth(str string) int64 {
	return int64(len(strings.Split(str, "/")) - 1)
}
