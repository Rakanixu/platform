package globals

import (
	"crypto/md5"
	"encoding/hex"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/slack"
)

const (
	NAMESPACE               string = "com.kazoup"
	FLAG_SERVICE_NAME       string = NAMESPACE + ".srv.flag"
	DB_SERVICE_NAME         string = NAMESPACE + ".srv.db"
	DATASOURCE_SERVICE_NAME string = NAMESPACE + ".srv.datasource"
	FilesTopic              string = NAMESPACE + ".topic.files"
	SlackChannelsTopic      string = NAMESPACE + ".topic.slackchannels"
	SlackUsersTopic         string = NAMESPACE + ".topic.slackusers"
	ScanTopic               string = NAMESPACE + ".topic.scan"
	CrawlerFinishedTopic    string = NAMESPACE + ".topic.crawlerfinished"
	NotificationTopic       string = NAMESPACE + ".topic.notification"

	IndexHelper = "files_helper"
	FilesAlias  = "files"
	FileType    = "file"

	FileTypeFile      = "files"
	FileTypeDirectory = "directories"

	Local       = "local"
	Slack       = "slack"
	GoogleDrive = "googledrive"
	OneDrive    = "onedrive"

	SlackFilesEndpoint    = "https://slack.com/api/files.list"
	SlackUsersEndpoint    = "https://slack.com/api/users.list"
	SlackChannelsEndpoint = "https://slack.com/api/channels.list"

	OneDriveEndpoint = "https://api.onedrive.com/v1.0/"

	GoogleDriveThumbnail = "https://www.googleapis.com/drive/v3/files/"

	OauthStateString = "randomsdsdahfoashfouahsfohasofhoashfaf"

	StartScanTask = "start_scan"
)

func NewGoogleOautConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/google/callback",
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
		RedirectURL:  "http://localhost:8082/auth/slack/callback",
		ClientID:     "2506087186.66729631906",
		ClientSecret: "53ea1f0afa4560b7e070964fb2b0c5d6",
		Scopes: []string{
			"files:read",
			"files:write:user",
			"team:read",
			"users:read",
			"channels:read",
		},
		Endpoint: slack.Endpoint,
	}
}

func NewMicrosoftOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  "http://localhost:8082/auth/microsoft/callback",
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

// Remove records (Files) from db that not longer belong to a datasource
// Compares LastSeen with the time the crawler started
// so all records with a LastSeen before will be removed from index
// file does not exists any more on datasource
func ClearIndex(e *datasource_proto.Endpoint) error {
	c := db_proto.NewDBClient(DB_SERVICE_NAME, nil)
	_, err := c.DeleteByQuery(context.Background(), &db_proto.DeleteByQueryRequest{
		Indexes:  []string{e.Index},
		Types:    []string{"file"},
		LastSeen: e.LastScanStarted,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
