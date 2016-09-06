package globals

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/slack"
)

const (
	FilesTopic           = "go.micro.topic.files"
	ScanTopic            = "go.micro.topic.scan"
	CrawlerFinishedTopic = "go.micro.topic.crawlerfinished"
	IndexHelper          = "files_helper"
	FilesAlias           = "files"
	FileType             = "file"

	Fake        = "fake"
	Local       = "local"
	Slack       = "slack"
	GoogleDrive = "googledrive"

	SlackFilesEndpoint = "https://slack.com/api/files.list"

	OauthStateString = "randomsdsdahfoashfouahsfohasofhoashfaf"
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
		Scopes:       []string{"files:read", "files:write:user", "team:read"},
		Endpoint:     slack.Endpoint,
	}
}
