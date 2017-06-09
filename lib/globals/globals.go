package globals

import (
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/slack"
	"time"
)

const (
	NAMESPACE                 string = "com.kazoup"
	DATASOURCE_SERVICE_NAME   string = NAMESPACE + ".srv.datasource"
	CRAWLER_SERVICE_NAME      string = NAMESPACE + ".srv.crawler"
	NOTIFICATION_SERVICE_NAME string = NAMESPACE + ".srv.notification"
	FILE_SERVICE_NAME         string = NAMESPACE + ".srv.file"
	USER_SERVICE_NAME         string = NAMESPACE + ".srv.user"
	CHANNEL_SERVICE_NAME      string = NAMESPACE + ".srv.channel"
	QUOTA_SERVICE_NAME        string = NAMESPACE + ".srv.quota"
	PROFILE_SERVICE_NAME      string = NAMESPACE + ".srv.profile"
	MONITOR_SERVICE_NAME      string = NAMESPACE + ".srv.monitor"
	THUMBNAIL_SERVICE_NAME    string = NAMESPACE + ".srv.thumbnail"
	AUDIO_SERVICE_NAME        string = NAMESPACE + ".srv.audio"
	DOCUMENT_SERVICE_NAME     string = NAMESPACE + ".srv.document"
	IMAGE_SERVICE_NAME        string = NAMESPACE + ".srv.image"
	ENTITIES_SERVICE_NAME     string = NAMESPACE + ".srv.entities"
	SENTIMENT_SERVICE_NAME    string = NAMESPACE + ".srv.sentiment"
	TRANSLATE_SERVICE_NAME    string = NAMESPACE + ".srv.translate"
    AGENT_SERVICE_NAME        string = NAMESPACE + ".srv.agent"

	HANDLER_DATASOURCE_CREATE  = DATASOURCE_SERVICE_NAME + ".Service.Create"
	HANDLER_DATASOURCE_DELETE  = DATASOURCE_SERVICE_NAME + ".Service.Delete"
	HANDLER_DATASOURCE_SEARCH  = DATASOURCE_SERVICE_NAME + ".Service.Search"
	HANDLER_DATASOURCE_SCAN    = DATASOURCE_SERVICE_NAME + ".Service.Scan"
	HANDLER_DATASOURCE_SCANALL = DATASOURCE_SERVICE_NAME + ".Service.ScanAll"

	HANDLER_FILE_CREATE = FILE_SERVICE_NAME + ".Service.Create"
	HANDLER_FILE_DELETE = FILE_SERVICE_NAME + ".Service.Delete"

	HANDLER_AUDIO_ENRICH_FILE       = AUDIO_SERVICE_NAME + ".Service.EnrichFile"
	HANDLER_AUDIO_ENRICH_DATASOURCE = AUDIO_SERVICE_NAME + ".Service.EnrichDatasource"

	HANDLER_DOCUMENT_ENRICH_FILE       = DOCUMENT_SERVICE_NAME + ".Service.EnrichFile"
	HANDLER_DOCUMENT_ENRICH_DATASOURCE = DOCUMENT_SERVICE_NAME + ".Service.EnrichDatasource"

	HANDLER_IMAGE_ENRICH_FILE       = IMAGE_SERVICE_NAME + ".Service.EnrichFile"
	HANDLER_IMAGE_ENRICH_DATASOURCE = IMAGE_SERVICE_NAME + ".Service.EnrichDatasource"

	HANDLER_SENTIMENT_ENRICH_FILE = SENTIMENT_SERVICE_NAME + ".Service.AnalyzeFile"

	HANDLER_ENTITIES_EXTRACT_FILE = ENTITIES_SERVICE_NAME + ".Service.ExtractFile"

	HANDLER_TRANSLATE_TRANSLATE = TRANSLATE_SERVICE_NAME + ".Service.Translate"
    HANDLER_AGENT_SAVE = AGENT_SERVICE_NAME + ".Service.Save"

	AnnounceTopic           string = NAMESPACE + ".topic.announce"
	DiscoverTopic           string = NAMESPACE + ".topic.discover.files"
	DiscoveryFinishedTopic  string = NAMESPACE + ".topic.discover.finish"
	DocEnrichTopic          string = NAMESPACE + ".topic.document"
	ImgEnrichTopic          string = NAMESPACE + ".topic.imgenrich"
	ThumbnailTopic          string = NAMESPACE + ".topic.thumbnail"
	AudioEnrichTopic        string = NAMESPACE + ".topic.audioenrich"
	SentimentEnrichTopic    string = NAMESPACE + ".topic.sentiment"
	ExtractEntitiesTopic    string = NAMESPACE + ".topic.extractentities"
	NotificationTopic       string = NAMESPACE + ".topic.notification"
	NotificationProxyTopic  string = NAMESPACE + ".topic.notificationproxy"
	DeleteBucketTopic       string = NAMESPACE + ".topic.deletebucket"
	DeleteFileInBucketTopic string = NAMESPACE + ".topic.deletefileinbucket"
    SaveRemoteFilesTopic        string = NAMESPACE + ".topic.saveremotefiles"

	IndexDatasources  = "datasources"
	IndexHelper       = "files_helper"
	FileType          = "file"
	FoldeType         = "folder"
	UserType          = "user"
	ChannelType       = "channel"
	FileTypeFile      = "files"
	FileTypeDirectory = "directories"
	TypeDatasource    = "datasource"

	Slack       = "slack"
	GoogleDrive = "googledrive"
	Gmail       = "gmail"
	OneDrive    = "onedrive"
	Dropbox     = "dropbox"
	Box         = "box"
	Mock        = "mock"

	DOCUMENT     = "document"
	PRESENTATION = "presentation"
	SPREADSHEET  = "spreadsheet"
	TEXT         = "text"

	SENSITIVE = "sensitive"

	MS_DOCUMENT     = "application/msword"
	MS_PRESENTATION = "application/vnd.ms-powerpoint"
	MS_SPREADSHEET  = "application/vnd.ms-excel"

	GOOGLE_DRIVE_DOCUMENT    = "application/vnd.google-apps.document"
	GOOGLE_DRIVE_PRESETATION = "application/vnd.google-apps.presentation"
	GOOGLE_DRIVE_SPREADSHEET = "application/vnd.google-apps.spreadsheet"
	GOOGLE_DRIVE_TEXT        = "text/plain"

	ONEDRIVE_TEXT = "text/plain"

	MIME_PNG  = "image/png"
	MIME_JPG  = "image/jpg"
	MIME_JPEG = "image/jpeg"

	CATEGORY_PICTURE  = "Pictures"
	CATEGORY_DOCUMENT = "Documents"
	CATEGORY_AUDIO    = "Audio"

	THUMBNAIL_WIDTH        = 178
	CLOUD_VISION_IMG_WIDTH = 640

	ACCESS_PUBLIC  = "public"
	ACCESS_SHARED  = "shared"
	ACCESS_PRIVATE = "private"

	GOOGLE_DRIVE_PUBLIC_FILE = "anyone"

	SlackFilesEndpoint       = "https://slack.com/api/files.list"
	SlackUsersEndpoint       = "https://slack.com/api/users.list"
	SlackChannelsEndpoint    = "https://slack.com/api/channels.list"
	SlackShareFilesEndpoint  = "https://slack.com/api/files.sharedPublicURL"
	SlackPostMessageEndpoint = "https://slack.com/api/chat.postMessage"

	OneDriveEndpoint = "https://api.onedrive.com/v1.0/"

	DropboxAccountEndpoint   = "https://api.dropboxapi.com/2/users/get_account"
	DropboxFilesEndpoint     = "https://api.dropboxapi.com/2/files/list_folder"
	DropboxFileEndpoint      = "https://api.dropboxapi.com/2/files/get_metadata"
	DropboxThumbnailEndpoint = "https://content.dropboxapi.com/2/files/get_thumbnail"
	DropboxFileMembers       = "https://api.dropboxapi.com/2/sharing/list_file_members"
	DropboxFileUpload        = "https://content.dropboxapi.com/2/files/upload"
	DropboxFileShare         = "https://api.dropboxapi.com/2/sharing/add_file_member"
	DropboxFileDelete        = "https://api.dropboxapi.com/2/files/delete"
	DropboxFileDownload      = "https://content.dropboxapi.com/2/files/download"
	DropboxSharedLinks       = "https://api.dropboxapi.com/2/sharing/list_shared_links"

	BoxAccountEndpoint      = "https://api.box.com/2.0/users/me"
	BoxFoldersEndpoint      = "https://api.box.com/2.0/folders/"
	BoxFileMetadataEndpoint = "https://api.box.com/2.0/files/"
	BoxUploadEndpoint       = "https://upload.box.com/api/2.0/files/content"

	GmailEndpoint = "https://mail.google.com/mail/u/"

	ROSETTE_ENTITIES_ENDPOINT  = "https://api.rosette.com/rest/v1/entities/"
	ROSETTE_SENTIMENT_ENDPOINT = "https://api.rosette.com/rest/v1/sentiment/"

	TMP_TOKEN_BUCKET = "tmp-token"
	AUDIO_BUCKET     = "kazoup-audio-bucket"

	StartScanTask = "start_scan"

	SERVER_ADDRESS        = "https://web.kazoup.io"
	SECURE_SERVER_ADDRESS = "https://web.kazoup.io:8082"

	SYSTEM_TOKEN      = "ajsdIgsnaloHFGis823jsdgyjTGDKijfcjk783JDUYFJyggvwejkxsnmbkjwpoj6483"
	DB_ACCESS_TOKEN   = "GSjsfduh3jskJHGuiU87y-skjaXXu7hpcMkdKghsojssio_98sushmpPpodvhakasdB"
	CLIENT_ID_SECRET  = "EC1FD9R5t6D3cs9CzPbgJaBJjshoVgrJrTs6U39scYzYF7HYyMlv_mal2IjLLaA9" // Auth0 RPC API client
	ENCRYTION_KEY_32  = "asjklasd766adfashj22kljasdhyfjkh"
	ROSETTE_API_KEY   = "c6872fa01aa45d59438d56831bb5b1a2"
	STRIPE_SECRET_KEY = "sk_test_udDr3n0RMjGr8vcwiCdNx3ao"
	STRIPE_PUBLIC_KEY = "pk_test_6z7qNSW5GZsLNyTz2hIrK0q5"

	NOTIFY_REFRESH_DATASOURCES = "refresh-datasources"
	NOTIFY_REFRESH_SEARCH      = "refresh-search"

	DISCOVERY_DELAY_MS = 10 * time.Millisecond
)

func NewGoogleOautConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/google/callback",
		ClientID:     "928848534435-kjubrqvl1sp50sfs3icemj2ma6v2an5j.apps.googleusercontent.com",
		ClientSecret: "zZAQz3zP5xnpLaA1S_q6YNhy",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.appdata",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/drive.metadata",
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			"https://www.googleapis.com/auth/drive.photos.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func NewSlackOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/slack/callback",
		ClientID:     "2506087186.66729631906",
		ClientSecret: "53ea1f0afa4560b7e070964fb2b0c5d6",
		Scopes: []string{
			"files:read",
			"files:write:user",
			"team:read",
			"users:read",
			"channels:read",
			"chat:write:user",
		},
		Endpoint: slack.Endpoint,
	}
}

func NewMicrosoftOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/microsoft/callback",
		ClientID:     "60f54c2b-6631-4bf4-ae45-01b5715cb881",
		ClientSecret: "COC67cMupbGdSCx1Omc3Z5g",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.live.com/oauth20_authorize.srf",
			TokenURL: "https://login.live.com/oauth20_token.srf",
		},
		Scopes: []string{
			"onedrive.readonly",
			"onedrive.readwrite",
			"onedrive.appfolder",
			"offline_access",
		},
	}
}

func NewDropboxOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL: SECURE_SERVER_ADDRESS + "/auth/dropbox/callback",
		//ClientID:     "6l5aj1fombrp6i7",
		ClientID: "882k4mhdmtza7y1",
		//ClientSecret: "nf8xar3qc1f32li",
		ClientSecret: "krhjkoim5u2a3v3",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropbox.com/1/oauth2/token",
		},
		Scopes: []string{},
	}
}

func NewBoxOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/box/callback",
		ClientID:     "8ryeu572aa5rk7iun56hsb0g7ta1oblp",
		ClientSecret: "An5sAtmY5KzlCvrAZgQ4rXQtBY3v6TwT",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://app.box.com/api/oauth2/authorize",
			TokenURL: "https://app.box.com/api/oauth2/token",
		},
		Scopes: []string{},
	}
}

func NewGmailOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/gmail/callback",
		ClientID:     "928848534435-kjubrqvl1sp50sfs3icemj2ma6v2an5j.apps.googleusercontent.com",
		ClientSecret: "zZAQz3zP5xnpLaA1S_q6YNhy",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/gmail.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func NewGoogleCloudOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SECURE_SERVER_ADDRESS + "/auth/google/callback",
		ClientID:     "928848534435-kjubrqvl1sp50sfs3icemj2ma6v2an5j.apps.googleusercontent.com",
		ClientSecret: "zZAQz3zP5xnpLaA1S_q6YNhy",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.appdata",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/drive.metadata",
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			"https://www.googleapis.com/auth/drive.photos.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

// NewSystemContext System context
func NewSystemContext() context.Context {
	return metadata.NewContext(context.TODO(), map[string]string{
		"Authorization": SYSTEM_TOKEN,
	})
}

// NewContextFromJWT
func NewContextFromJWT(jwt string) context.Context {
	return metadata.NewContext(context.TODO(), map[string]string{
		"Authorization": jwt,
	})
}

var FileTypeDict = struct {
	M map[string]map[string]string
}{
	M: map[string]map[string]string{
		GoogleDrive: map[string]string{
			DOCUMENT:     GOOGLE_DRIVE_DOCUMENT,
			PRESENTATION: GOOGLE_DRIVE_PRESETATION,
			SPREADSHEET:  GOOGLE_DRIVE_SPREADSHEET,
			TEXT:         GOOGLE_DRIVE_TEXT,
		},
		OneDrive: map[string]string{
			TEXT: ONEDRIVE_TEXT,
		},
		Gmail: map[string]string{ //TODO: everything exept gdrive, but model is done
			GOOGLE_DRIVE_DOCUMENT:    "application/vnd.google-apps.document",
			GOOGLE_DRIVE_PRESETATION: "application/vnd.google-apps.presentation",
			GOOGLE_DRIVE_SPREADSHEET: "application/vnd.google-apps.spreadsheet",
			GOOGLE_DRIVE_TEXT:        "application/vnd.google-apps.file",
		},
	},
}
