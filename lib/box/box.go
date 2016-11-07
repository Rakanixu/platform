package box

type BoxDirContents struct {
	Type           string `json:"type"`
	ID             string `json:"id"`
	SequenceID     string `json:"sequence_id"`
	Etag           string `json:"etag"`
	Name           string `json:"name"`
	CreatedAt      string `json:"created_at"`
	ModifiedAt     string `json:"modified_at"`
	Description    string `json:"description"`
	Size           int    `json:"size"`
	PathCollection struct {
		TotalCount int `json:"total_count"`
		Entries    []struct {
			Type       string      `json:"type"`
			ID         string      `json:"id"`
			SequenceID interface{} `json:"sequence_id"`
			Etag       interface{} `json:"etag"`
			Name       string      `json:"name"`
		} `json:"entries"`
	} `json:"path_collection"`
	CreatedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"created_by"`
	ModifiedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"modified_by"`
	OwnedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"owned_by"`
	SharedLink struct {
		URL               string      `json:"url"`
		DownloadURL       interface{} `json:"download_url"`
		VanityURL         interface{} `json:"vanity_url"`
		IsPasswordEnabled bool        `json:"is_password_enabled"`
		UnsharedAt        interface{} `json:"unshared_at"`
		DownloadCount     int         `json:"download_count"`
		PreviewCount      int         `json:"preview_count"`
		Access            string      `json:"access"`
		Permissions       struct {
			CanDownload bool `json:"can_download"`
			CanPreview  bool `json:"can_preview"`
		} `json:"permissions"`
	} `json:"shared_link"`
	FolderUploadEmail struct {
		Access string `json:"access"`
		Email  string `json:"email"`
	} `json:"folder_upload_email"`
	Parent struct {
		Type       string      `json:"type"`
		ID         string      `json:"id"`
		SequenceID interface{} `json:"sequence_id"`
		Etag       interface{} `json:"etag"`
		Name       string      `json:"name"`
	} `json:"parent"`
	ItemStatus     string `json:"item_status"`
	ItemCollection struct {
		TotalCount int       `json:"total_count"`
		Entries    []BoxFile `json:"entries"`
		Offset     int       `json:"offset"`
		Limit      int       `json:"limit"`
	} `json:"item_collection"`
}

type BoxFileMeta struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	FileVersion struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Sha1 string `json:"sha1"`
	} `json:"file_version"`
	SequenceID     string `json:"sequence_id"`
	Etag           string `json:"etag"`
	Sha1           string `json:"sha1"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Size           int    `json:"size"`
	PathCollection struct {
		TotalCount int `json:"total_count"`
		Entries    []struct {
			Type       string      `json:"type"`
			ID         string      `json:"id"`
			SequenceID interface{} `json:"sequence_id"`
			Etag       interface{} `json:"etag"`
			Name       string      `json:"name"`
		} `json:"entries"`
	} `json:"path_collection"`
	CreatedAt         string      `json:"created_at"`
	ModifiedAt        string      `json:"modified_at"`
	TrashedAt         interface{} `json:"trashed_at"`
	PurgedAt          interface{} `json:"purged_at"`
	ContentCreatedAt  string      `json:"content_created_at"`
	ContentModifiedAt string      `json:"content_modified_at"`
	CreatedBy         struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"created_by"`
	ModifiedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"modified_by"`
	OwnedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Login string `json:"login"`
	} `json:"owned_by"`
	SharedLink struct {
		Access            string `json:"access"`
		DownloadCount     int    `json:"download_count"`
		DownloadURL       string `json:"download_url"`
		EffectiveAccess   string `json:"effective_access"`
		IsPasswordEnabled bool   `json:"is_password_enabled"`
		Permissions       struct {
			CanDownload bool `json:"can_download"`
			CanPreview  bool `json:"can_preview"`
		} `json:"permissions"`
		PreviewCount int         `json:"preview_count"`
		UnsharedAt   interface{} `json:"unshared_at"`
		URL          string      `json:"url"`
		VanityURL    interface{} `json:"vanity_url"`
	} `json:"shared_link"`
	Parent struct {
		Type       string      `json:"type"`
		ID         string      `json:"id"`
		SequenceID interface{} `json:"sequence_id"`
		Etag       interface{} `json:"etag"`
		Name       string      `json:"name"`
	} `json:"parent"`
	ItemStatus string `json:"item_status"`
}

type BoxFile struct {
	Type        string `json:"type"`
	ID          string `json:"id"`
	FileVersion struct {
		Type string `json:"type"`
		ID   string `json:"id"`
		Sha1 string `json:"sha1"`
	} `json:"file_version"`
	SequenceID string `json:"sequence_id"`
	Etag       string `json:"etag"`
	Sha1       string `json:"sha1"`
	Name       string `json:"name"`
}

type BoxUpload struct {
	TotalCount int           `json:"total_count"`
	Entries    []BoxFileMeta `json:"entries"`
}
