package onedrive

import "time"

type DrivesListResponse struct {
	OdataContext string `json:"@odata.context,omitempty"`
	Value        []struct {
		ID        string `json:"id,omitempty"`
		DriveType string `json:"driveType,omitempty"`
		Owner     struct {
			User struct {
				DisplayName string `json:"displayName,omitempty"`
				ID          string `json:"id,omitempty"`
			} `json:"user,omitempty"`
		} `json:"owner,omitempty"`
		Quota struct {
			Deleted      int    `json:"deleted,omitempty"`
			Remaining    int64  `json:"remaining,omitempty"`
			State        string `json:"state,omitempty"`
			Total        int64  `json:"total,omitempty"`
			Used         int    `json:"used,omitempty"`
			StoragePlans struct {
				UpgradeAvailable bool `json:"upgradeAvailable,omitempty"`
			} `json:"storagePlans,omitempty"`
		} `json:"quota,omitempty"`
		Status struct {
			State string `json:"state,omitempty"`
		} `json:"status,omitempty"`
	} `json:"value,omitempty"`
}

type FilesListResponse struct {
	OdataContext string `json:"@odata.context,omitempty"`
	Value        []OneDriveFile
}

type OneDriveFile struct {
	ContentDownloadURL string `json:"contentDdownloadUrl,omitempty"`
	CreatedBy          struct {
		Application struct {
			DisplayName string `json:"displayName,omitempty"`
			ID          string `json:"id,omitempty"`
		} `json:"application,omitempty"`
		User struct {
			DisplayName string `json:"displayName,omitempty"`
			ID          string `json:"id,omitempty"`
		} `json:"user,omitempty"`
	} `json:"createdBy,omitempty"`
	CreatedDateTime time.Time `json:"createdDateTime,omitempty"`
	CTag            string    `json:"cTag,omitempty"`
	ETag            string    `json:"eTag,omitempty"`
	ID              string    `json:"id,omitempty"`
	LastModifiedBy  struct {
		Application struct {
			DisplayName string `json:"displayName,omitempty"`
			ID          string `json:"id,omitempty"`
		} `json:"application,omitempty"`
		User struct {
			DisplayName string `json:"displayName,omitempty"`
			ID          string `json:"id,omitempty"`
		} `json:"user,omitempty"`
	} `json:"lastModifiedBy,omitempty"`
	LastModifiedDateTime time.Time `json:"lastModifiedDateTime,omitempty"`
	Name                 string    `json:"name,omitempty"`
	ParentReference      struct {
		DriveID string `json:"driveId,omitempty"`
		ID      string `json:"id,omitempty"`
		Path    string `json:"path,omitempty"`
	} `json:"parentReference,omitempty"`
	Size           int64  `json:"size,omitempty"`
	WebURL         string `json:"webUrl,omitempty"`
	FileSystemInfo struct {
		CreatedDateTime      time.Time `json:"createdDateTime,omitempty"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime,omitempty"`
	} `json:"fileSystemInfo,omitempty"`
	Folder struct {
		ChildCount int `json:"childCount,omitempty"`
	} `json:"folder,omitempty"`
	SpecialFolder struct {
		Name string `json:"name,omitempty"`
	} `json:"specialFolder,omitempty"`
	File struct {
		Hashes struct {
			Crc32Hash string `json:"crc32Hash,omitempty"`
			Sha1Hash  string `json:"sha1Hash,omitempty"`
		} `json:"hashes,omitempty"`
		MimeType string `json:"mimeType,omitempty"`
	} `json:"file,omitempty"`
}
