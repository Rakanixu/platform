package slack

// FilesListResponse represents https://slack.com/api/files.list response
type FilesListResponse struct {
	Ok     bool `json:"ok"`
	Files  []SlackFile
	Paging Page `json:"paging"`
}

type Page struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Pages int `json:"pages"`
}

// SlackFile
type SlackFile struct {
	ID                 string        `json:"id"`
	Created            int           `json:"created"`
	Timestamp          int64         `json:"timestamp"`
	Name               string        `json:"name"`
	Title              string        `json:"title"`
	Mimetype           string        `json:"mimetype"`
	Filetype           string        `json:"filetype"`
	PrettyType         string        `json:"pretty_type"`
	User               string        `json:"user"`
	Editable           bool          `json:"editable"`
	Size               int64         `json:"size"`
	Mode               string        `json:"mode"`
	IsExternal         bool          `json:"is_external"`
	ExternalType       string        `json:"external_type"`
	IsPublic           bool          `json:"is_public"`
	PublicURLShared    bool          `json:"public_url_shared"`
	DisplayAsBot       bool          `json:"display_as_bot"`
	Username           string        `json:"username"`
	URLPrivate         string        `json:"url_private"`
	URLPrivateDownload string        `json:"url_private_download"`
	Thumb64            string        `json:"thumb_64"`
	Thumb80            string        `json:"thumb_80"`
	Thumb360           string        `json:"thumb_360"`
	Thumb360W          int           `json:"thumb_360_w"`
	Thumb360H          int           `json:"thumb_360_h"`
	Thumb480           string        `json:"thumb_480"`
	Thumb480W          int           `json:"thumb_480_w"`
	Thumb480H          int           `json:"thumb_480_h"`
	Thumb160           string        `json:"thumb_160"`
	Thumb720           string        `json:"thumb_720"`
	Thumb720W          int           `json:"thumb_720_w"`
	Thumb720H          int           `json:"thumb_720_h"`
	Thumb960           string        `json:"thumb_960"`
	Thumb960W          int           `json:"thumb_960_w"`
	Thumb960H          int           `json:"thumb_960_h"`
	Thumb1024          string        `json:"thumb_1024"`
	Thumb1024W         int           `json:"thumb_1024_w"`
	Thumb1024H         int           `json:"thumb_1024_h"`
	ImageExifRotation  int           `json:"image_exif_rotation"`
	OriginalW          int           `json:"original_w"`
	OriginalH          int           `json:"original_h"`
	Permalink          string        `json:"permalink"`
	PermalinkPublic    string        `json:"permalink_public"`
	Channels           []string      `json:"channels"`
	Groups             []interface{} `json:"groups"`
	Ims                []interface{} `json:"ims"`
	CommentsCount      int           `json:"comments_count"`
	InitialComment     struct {
		ID        string `json:"id"`
		Created   int    `json:"created"`
		Timestamp int    `json:"timestamp"`
		User      string `json:"user"`
		IsIntro   bool   `json:"is_intro"`
		Comment   string `json:"comment"`
		Channel   string `json:"channel"`
	} `json:"initial_comment"`
}

type ChannelListResponse struct {
	Ok       bool `json:"ok"`
	Channels []SlackChannel
}

type ESSlackChannel struct {
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	Created     int    `json:"created"`
	Creator     string `json:"creator"`
	IsArchived  bool   `json:"is_archived"`
	IsMember    bool   `json:"is_member"`
	NumMembers  int    `json:"num_members"`
	Topic       struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"topic"`
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"purpose"`
}

type SlackChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	ESSlackChannel
}

type UserListResponse struct {
	Ok      bool        `json:"ok"`
	Members []SlackUser `json:"members"`
}

type ESSlackUser struct {
	UserID   string      `json:"user_id"`
	TeamID   string      `json:"team_id"`
	UserName string      `json:"user_name"`
	Deleted  bool        `json:"deleted"`
	Status   interface{} `json:"status,omitempty"`
	Color    string      `json:"color,omitempty"`
	RealName string      `json:"real_name,omitempty"`
	Tz       interface{} `json:"tz,omitempty"`
	TzLabel  string      `json:"tz_label,omitempty"`
	TzOffset int         `json:"tz_offset,omitempty"`
	Profile  struct {
		BotID              string `json:"bot_id"`
		APIAppID           string `json:"api_app_id"`
		FirstName          string `json:"first_name"`
		AvatarHash         string `json:"avatar_hash"`
		Image24            string `json:"image_24"`
		Image32            string `json:"image_32"`
		Image48            string `json:"image_48"`
		Image72            string `json:"image_72"`
		Image192           string `json:"image_192"`
		Image512           string `json:"image_512"`
		Image1024          string `json:"image_1024"`
		ImageOriginal      string `json:"image_original"`
		RealName           string `json:"real_name"`
		RealNameNormalized string `json:"real_name_normalized"`
	} `json:"profile"`
	IsAdmin           bool   `json:"is_admin,omitempty"`
	IsOwner           bool   `json:"is_owner,omitempty"`
	IsPrimaryOwner    bool   `json:"is_primary_owner,omitempty"`
	IsRestricted      bool   `json:"is_restricted,omitempty"`
	IsUltraRestricted bool   `json:"is_ultra_restricted,omitempty"`
	IsBot             bool   `json:"is_bot,omitempty"`
	Has2Fa            bool   `json:"has_2fa,omitempty"`
	TwoFactorType     string `json:"two_factor_type,omitempty"`
}

type SlackUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	ESSlackUser
}

type SlackShareResponse struct {
	Ok   bool      `json:"ok"`
	File SlackFile `json:"file"`
}
