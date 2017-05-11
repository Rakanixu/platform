package file

import (
	"encoding/json"
	"fmt"
	"github.com/kazoup/platform/lib/box"
	"github.com/kazoup/platform/lib/categories"
	"github.com/kazoup/platform/lib/dropbox"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	gmailhelper "github.com/kazoup/platform/lib/gmail"
	"github.com/kazoup/platform/lib/onedrive"
	"github.com/kazoup/platform/lib/slack"
	googledrive "google.golang.org/api/drive/v3"
	"net/url"
	"strings"
	"time"
)

func NewFileFromString(s string) (File, error) {
	kf := &KazoupFile{}
	if err := json.Unmarshal([]byte(s), kf); err != nil {
		return nil, err
	}

	switch kf.FileType {
	case globals.Slack:
		ksf := &KazoupSlackFile{}
		if err := json.Unmarshal([]byte(s), ksf); err != nil {
			return nil, err
		}
		return ksf, nil
	case globals.GoogleDrive:
		kgf := &KazoupGoogleFile{}
		if err := json.Unmarshal([]byte(s), kgf); err != nil {
			return nil, err
		}
		return kgf, nil
	case globals.Gmail:
		kgf := &KazoupGmailFile{}
		if err := json.Unmarshal([]byte(s), kgf); err != nil {
			return nil, err
		}
		return kgf, nil
	case globals.OneDrive:
		kof := &KazoupOneDriveFile{}
		if err := json.Unmarshal([]byte(s), kof); err != nil {
			return nil, err
		}
		return kof, nil
	case globals.Dropbox:
		kdf := &KazoupDropboxFile{}
		if err := json.Unmarshal([]byte(s), kdf); err != nil {
			return nil, err
		}
		return kdf, nil
	case globals.Box:
		kbf := &KazoupBoxFile{}
		if err := json.Unmarshal([]byte(s), kbf); err != nil {
			return nil, err
		}
		return kbf, nil
	case globals.Mock:
		kmf := &KazoupMockFile{}
		if err := json.Unmarshal([]byte(s), kmf); err != nil {
			return nil, err
		}
		return kmf, nil
	default:
		return nil, errors.ErrInvalidFile
	}
}

// NewKazoupFileFromGoogleDriveFile constructor
// opts first param is base64 emcoded file (itself)
func NewKazoupFileFromGoogleDriveFile(g googledrive.File, dsId, uId, index string) *KazoupGoogleFile {
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", g.ModifiedTime)
	access := globals.ACCESS_PRIVATE
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

	// More permissions than owner
	if len(g.Permissions) > 1 {
		access = globals.ACCESS_SHARED
	}

	// We will overwrite access if we find a type anyone as permission
	for _, v := range g.Permissions {
		if v.Type == globals.GOOGLE_DRIVE_PUBLIC_FILE {
			access = globals.ACCESS_PUBLIC
			break
		}
	}

	// Do not index trashed files
	if g.Trashed {
		return nil
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
		MimeType:     g.MimeType,
		Depth:        0,
		FileType:     globals.GoogleDrive,
		LastSeen:     time.Now().Unix(),
		Access:       access,
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupGoogleFile{*kf, &g}
}

// NewKazoupFileFromSlackFile constructor
func NewKazoupFileFromSlackFile(s slack.SlackFile, dsId, uId, index string) *KazoupSlackFile {
	t := time.Unix(s.Timestamp, 0)
	access := globals.ACCESS_PRIVATE

	if s.IsPublic { // Is public to the the team, not globally
		access = globals.ACCESS_SHARED
	}

	if s.PublicURLShared && s.Mode == "hosted" {
		access = globals.ACCESS_PUBLIC
	}

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(s.URLPrivate),
		UserId:       uId,
		Name:         s.Name,
		URL:          s.URLPrivate,
		Modified:     t,
		FileSize:     s.Size,
		IsDir:        false,
		Category:     categories.GetDocType("." + s.Filetype),
		MimeType:     s.Mimetype,
		Depth:        0,
		FileType:     globals.Slack,
		LastSeen:     time.Now().Unix(),
		Access:       access,
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupSlackFile{*kf, &s}
}

// NewKazoupFileFromOneDriveFile constructor
func NewKazoupFileFromOneDriveFile(o onedrive.OneDriveFile, dsId, uId, index string) *KazoupOneDriveFile {
	isDir := true
	mimeType := ""
	access := globals.ACCESS_PRIVATE
	name := strings.Split(o.Name, ".")

	if len(o.File.MimeType) > 0 {
		isDir = false
		mimeType = o.File.MimeType
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
		MimeType:     mimeType,
		Depth:        0,
		FileType:     globals.OneDrive,
		LastSeen:     time.Now().Unix(),
		Access:       access,
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupOneDriveFile{*kf, &o}
}

// NewKazoupFileFromDropboxFile constructor
func NewKazoupFileFromDropboxFile(d dropbox.DropboxFile, dsId, uId, index string) *KazoupDropboxFile {
	isDir := false
	name := strings.Split(d.Name, ".")
	path := strings.Replace(d.PathDisplay, "/"+d.Name, "", 1)
	url := fmt.Sprintf("https://www.dropbox.com/home%s?preview=%s", path, url.QueryEscape(d.Name))

	// Dropbox file fall into those categories: file, folder, deleted
	// By default, deleted files AND deleted folders will be flag as (isDir = false), then will appear on the frontend
	if d.Tag == globals.FoldeType {
		isDir = true
	}

	d.DropboxTag = d.Tag // Store the tag friendly

	// Do not index trashed files, as in dropbox generate rubish for basic account.
	// Trashed files are not longer there (there are but flaged in dropbox, probably an upgrade account make them reachable)
	if d.DropboxTag == "deleted" {
		return nil
	}

	kf := &KazoupFile{
		ID:           globals.GetMD5Hash(url),
		UserId:       uId,
		Name:         d.Name,
		URL:          url,
		Modified:     d.ServerModified,
		FileSize:     int64(d.Size),
		IsDir:        isDir,
		Category:     categories.GetDocType("." + name[len(name)-1]),
		MimeType:     "", // TODO: mime type Not present on origial file
		Depth:        0,
		FileType:     globals.Dropbox,
		LastSeen:     time.Now().Unix(),
		Access:       "", // Filled after calling this constructor on the fs.Walk
		DatasourceId: dsId,
		Index:        index,
	}

	return &KazoupDropboxFile{*kf, &d}
}

// NewKazoupFileFromDropboxFile constructor
func NewKazoupFileFromBoxFile(d box.BoxFileMeta, dsId, uId, index string) *KazoupBoxFile {
	access := globals.ACCESS_PRIVATE
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

	if len(d.SharedLink.URL) > 0 {
		access = globals.ACCESS_PUBLIC
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
		MimeType:     "", // TODO: mime type Not present on origial file
		Depth:        0,
		FileType:     globals.Box,
		LastSeen:     time.Now().Unix(),
		Access:       access,
		DatasourceId: dsId,
		Index:        index,
	}

	return &KazoupBoxFile{*kf, &d}
}

// NewKazoupFileFromGmailFile constructor
func NewKazoupFileFromGmailFile(m gmailhelper.GmailFile, dsId, uId, dsURL, index string) *KazoupGmailFile {
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
		MimeType:     m.MimeType,
		Depth:        0,
		FileType:     globals.Gmail,
		LastSeen:     time.Now().Unix(),
		Access:       globals.ACCESS_PRIVATE, // Mail attachments are always private, once send, original owner cannot delete that file for other user
		DatasourceId: dsId,
		Index:        index,
	}
	return &KazoupGmailFile{*kf, &m}
}

// NewKazoupFileFromMockFile constructor
func NewKazoupFileFromMockFile() *KazoupMockFile {
	return &KazoupMockFile{KazoupFile: KazoupFile{
		ID:           globals.Mock,
		UserId:       globals.Mock,
		Name:         globals.Mock,
		FileType:     globals.Mock,
		DatasourceId: globals.Mock,
		Index:        globals.Mock,
	}}
}
