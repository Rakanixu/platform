package file

import (
	"encoding/json"
	"errors"
	"fmt"
	db "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/local"
	"github.com/kazoup/platform/lib/onedrive"
	"github.com/kazoup/platform/lib/slack"
	"golang.org/x/net/context"
	googledrive "google.golang.org/api/drive/v3"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

// GetFileByID retrieves a file given its id and the user belongs to
func GetFileByID(ctx context.Context, md5UserId, id string, c db.DBClient) (File, error) {
	dbres, err := c.SearchById(ctx, &db.SearchByIdRequest{
		Index: md5UserId,
		Type:  "file",
		Id:    id,
	})
	if err != nil {
		return nil, err
	}

	f, err := NewFileFromString(dbres.Result)
	if err != nil {
		return nil, err
	}

	return f, err
}

func NewFileFromString(s string) (File, error) {
	kf := &KazoupFile{}
	if err := json.Unmarshal([]byte(s), kf); err != nil {
		return nil, errors.New("Error unmarsahling NewFileFromString")
	}

	switch kf.FileType {
	case globals.Local:
		klf := &KazoupLocalFile{}
		if err := json.Unmarshal([]byte(s), klf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case local")
		}
		return klf, nil
	case globals.Slack:
		ksf := &KazoupSlackFile{}
		if err := json.Unmarshal([]byte(s), ksf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case slack")
		}
		return ksf, nil
	case globals.GoogleDrive:
		kgf := &KazoupGoogleFile{}
		if err := json.Unmarshal([]byte(s), kgf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case googledrive")
		}
		return kgf, nil
	case globals.Gmail:
		kgf := &KazoupGmailFile{}
		if err := json.Unmarshal([]byte(s), kgf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case gmail")
		}
		return kgf, nil
	case globals.OneDrive:
		kof := &KazoupOneDriveFile{}
		if err := json.Unmarshal([]byte(s), kof); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case onedrive")
		}
		return kof, nil
	case globals.Dropbox:
		kdf := &KazoupDropboxFile{}
		if err := json.Unmarshal([]byte(s), kdf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case dropbox")
		}
		return kdf, nil
	case globals.Box:
		kbf := &KazoupBoxFile{}
		if err := json.Unmarshal([]byte(s), kbf); err != nil {
			return nil, errors.New("Error unmarsahling NewFileFromString case box")
		}
		return kbf, nil
	default:
		return nil, errors.New("Error constructing file type")
	}
}

// NewKazoupFileFromGoogleDriveFile constructor
func NewKazoupFileFromGoogleDriveFile(g *googledrive.File, dsId, uId, index string) *KazoupGoogleFile {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", g.ModifiedTime)

	d := false
	if len(g.FolderColorRgb) > 0 {
		d = true
	}

	c := categories.GetDocType("." + g.FullFileExtension)
	if len(g.FullFileExtension) == 0 {
		c = categories.GetDocType(g.MimeType)
	}

	// Link to view / edit file
	url := g.WebViewLink
	// When document is Microsoft Office, download link
	if g.MimeType == globals.MS_DOCUMENT ||
		g.MimeType == globals.MS_PRESENTATION ||
		g.MimeType == globals.MS_SPREADSHEET {
		url = g.WebContentLink
	}

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(url),
		UserId:       uId,
		Name:         g.Name,
		URL:          url,
		Modified:     t,
		FileSize:     g.Size,
		IsDir:        d,
		Category:     c,
		Depth:        0,
		FileType:     globals.GoogleDrive,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupGoogleFile{*kf, *g}
}

// NewKazoupFileFromSlackFile constructor
func NewKazoupFileFromSlackFile(s *slack.SlackFile, dsId, uId, index string) *KazoupSlackFile {
	t := time.Unix(s.Timestamp, 0)

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(s.URLPrivate),
		UserId:       uId,
		Name:         s.Name,
		URL:          s.URLPrivate,
		Modified:     t,
		FileSize:     s.Size,
		IsDir:        false,
		Category:     categories.GetDocType("." + s.Filetype),
		Depth:        0,
		FileType:     globals.Slack,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupSlackFile{*kf, *s}
}

func NewKazoupFileFromLocal(lf *local.LocalFile, dsId, uId, index string) *KazoupLocalFile {
	// don;t save all LocalFile as mmost of data is same as KazoupFile just pass file mode
	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(lf.Path),
		UserId:       uId,
		Name:         lf.Info.Name(),
		URL:          "/local" + lf.Path,
		Modified:     lf.Info.ModTime(),
		FileSize:     lf.Info.Size(),
		IsDir:        lf.Info.IsDir(),
		Category:     categories.GetDocType(filepath.Ext(lf.Info.Name())),
		Depth:        UrlDepth(lf.Path),
		FileType:     globals.Local,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupLocalFile{*kf}

}

// NewKazoupFileFromOneDriveFile constructor
func NewKazoupFileFromOneDriveFile(o *onedrive.OneDriveFile, dsId, uId, index string) *KazoupOneDriveFile {
	isDir := true
	name := strings.Split(o.Name, ".")

	if len(o.File.MimeType) > 0 {
		isDir = false
	}

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(o.WebURL),
		UserId:       uId,
		Name:         o.Name,
		URL:          o.WebURL,
		Modified:     o.LastModifiedDateTime,
		FileSize:     o.Size,
		IsDir:        isDir,
		Category:     categories.GetDocType("." + name[len(name)-1]),
		Depth:        0,
		FileType:     globals.OneDrive,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupOneDriveFile{*kf, *o}
}

// NewKazoupFileFromDropboxFile constructor
func NewKazoupFileFromDropboxFile(d *dropbox.DropboxFile, dsId, uId, index string) *KazoupDropboxFile {
	isDir := false
	name := strings.Split(d.Name, ".")
	path := strings.Replace(d.PathDisplay, "/"+d.Name, "", 1)
	url := fmt.Sprintf("https://www.dropbox.com/home%s?preview=%s", path, url.QueryEscape(d.Name))

	// Dropbox file fall into those categories: file, folder, deleted
	// By default, deleted files AND deleted folders will be flag as (isDir = false), then will appear on the frontend
	// On the frontend
	if d.Tag == globals.FoldeType {
		isDir = true
	}

	d.DropboxTag = d.Tag // Store the tag friendly

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(url),
		UserId:       uId,
		Name:         d.Name,
		URL:          url,
		Modified:     d.ServerModified,
		FileSize:     int64(d.Size),
		IsDir:        isDir,
		Category:     categories.GetDocType("." + name[len(name)-1]),
		Depth:        0,
		FileType:     globals.Dropbox,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}

	return &KazoupDropboxFile{*kf, *d}
}

// NewKazoupFileFromDropboxFile constructor
func NewKazoupFileFromBoxFile(d *box.BoxFileMeta, dsId, uId, index string) *KazoupBoxFile {
	isDir := false
	name := strings.Split(d.Name, ".")
	url := fmt.Sprintf("https://app.box.com/%s/%s", d.Type, d.ID)
	t, err := time.Parse(time.RFC3339, d.ModifiedAt)
	if err != nil {
		fmt.Println(err)
	}

	if d.Type == "folder" {
		isDir = true
	}

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(url),
		UserId:       uId,
		Name:         d.Name,
		URL:          url,
		Modified:     t,
		FileSize:     int64(d.Size),
		IsDir:        isDir,
		Category:     categories.GetDocType("." + name[len(name)-1]),
		Depth:        0,
		FileType:     globals.Box,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}

	return &KazoupBoxFile{*kf, *d}
}

// NewKazoupFileFromGmailFile constructor
func NewKazoupFileFromGmailFile(m *gmailhelper.GmailFile, dsId, uId, dsURL, index string) *KazoupGmailFile {
	if len(m.Name) == 0 {
		return nil
	}

	url := fmt.Sprintf("%s%s/#inbox/%s", globals.GmailEndpoint, strings.Replace(dsURL, globals.Gmail+"://", "", 1), m.MessageId)
	t := time.Unix(m.InternalDate/1000, 0)

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(url),
		UserId:       uId,
		Name:         m.Name,
		URL:          url,
		Modified:     t,
		FileSize:     m.SizeEstimate,
		IsDir:        false,
		Category:     categories.GetDocType(fmt.Sprintf(".%s", m.Extension)),
		Depth:        0,
		FileType:     globals.Gmail,
		LastSeen:     time.Now().Unix(),
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupGmailFile{*kf, *m}
}

func UrlDepth(str string) int64 {
	return int64(len(strings.Split(str, "/")) - 1)
}
