package onedrive

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/kazoup/platform/structs"
	"github.com/kazoup/platform/structs/categories"
	"os"
	"strings"
	"time"
)

type DrivesListResponse struct {
	OdataContext string `json:"@odata.context"`
	Value        []struct {
		ID        string `json:"id"`
		DriveType string `json:"driveType"`
		Owner     struct {
			User struct {
				DisplayName string `json:"displayName"`
				ID          string `json:"id"`
			} `json:"user"`
		} `json:"owner"`
		Quota struct {
			Deleted      int    `json:"deleted"`
			Remaining    int64  `json:"remaining"`
			State        string `json:"state"`
			Total        int64  `json:"total"`
			Used         int    `json:"used"`
			StoragePlans struct {
				UpgradeAvailable bool `json:"upgradeAvailable"`
			} `json:"storagePlans"`
		} `json:"quota"`
		Status struct {
			State string `json:"state"`
		} `json:"status"`
	} `json:"value"`
}

type FilesListResponse struct {
	OdataContext string `json:"@odata.context"`
	Value        []OneDriveFile
}

type OneDriveFile struct {
	ContentDownloadURL string `json:"contentDdownloadUrl"`
	CreatedBy          struct {
		Application struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"application"`
		User struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"user"`
	} `json:"createdBy"`
	CreatedDateTime time.Time `json:"createdDateTime"`
	CTag            string    `json:"cTag"`
	ETag            string    `json:"eTag"`
	ID              string    `json:"id"`
	LastModifiedBy  struct {
		Application struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"application"`
		User struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
		} `json:"user"`
	} `json:"lastModifiedBy"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	Name                 string    `json:"name"`
	ParentReference      struct {
		DriveID string `json:"driveId"`
		ID      string `json:"id"`
		Path    string `json:"path"`
	} `json:"parentReference"`
	Size           int64  `json:"size"`
	WebURL         string `json:"webUrl"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime"`
	} `json:"fileSystemInfo"`
	Folder struct {
		ChildCount int `json:"childCount"`
	} `json:"folder"`
	SpecialFolder struct {
		Name string `json:"name"`
	} `json:"specialFolder"`
	File struct {
		Hashes struct {
			Crc32Hash string `json:"crc32Hash"`
			Sha1Hash  string `json:"sha1Hash"`
		} `json:"hashes"`
		MimeType string `json:"mimeType"`
	} `json:"file"`
}

// NewKazoupFileFromOneDriveFile constructor
func NewKazoupFileFromOneDriveFile(o *OneDriveFile) *structs.KazoupFile {
	isDir := true
	name := strings.Split(o.Name, ".")

	if len(o.File.MimeType) > 0 {
		isDir = false
	}

	return &structs.KazoupFile{
		ID:       getMD5Hash(o.WebURL),
		Name:     o.Name,
		URL:      o.WebURL,
		Modified: o.LastModifiedDateTime,
		Size:     o.Size,
		IsDir:    isDir,
		Mode:     *new(os.FileMode),
		Category: categories.GetDocType("." + name[len(name)-1]),
		Depth:    0,
		Original: o,
	}
}

func getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
