package globals

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/slack"
	"io"
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

	IndexDatasources = "datasources"
	IndexFlags       = "flags"
	IndexHelper      = "files_helper"
	FilesAlias       = "files"
	FileType         = "file"

	FileTypeFile      = "files"
	FileTypeDirectory = "directories"

	Local       = "local"
	Slack       = "slack"
	GoogleDrive = "googledrive"
	Gmail       = "gmail"
	OneDrive    = "onedrive"
	Dropbox     = "dropbox"
	Box         = "box"

	SlackFilesEndpoint    = "https://slack.com/api/files.list"
	SlackUsersEndpoint    = "https://slack.com/api/users.list"
	SlackChannelsEndpoint = "https://slack.com/api/channels.list"

	OneDriveEndpoint = "https://api.onedrive.com/v1.0/"

	DropboxAccountEndpoint   = "https://api.dropboxapi.com/2/users/get_account"
	DropboxFilesEndpoint     = "https://api.dropboxapi.com/2/files/list_folder"
	DropboxThumbnailEndpoint = "https://content.dropboxapi.com/2/files/get_thumbnail"

	BoxAccountEndpoint      = "https://api.box.com/2.0/users/me"
	BoxFoldersEndpoint      = "https://api.box.com/2.0/folders/"
	BoxFileMetadataEndpoint = "https://api.box.com/2.0/files/"

	GmailEndpoint = "https://mail.google.com/mail/u/"

	StartScanTask = "start_scan"

	SERVER_ADDRESS        = "http://web.kazoup.io:8082"
	SECURE_SERVER_ADDRESS = "https://web.kazoup.io:8082"

	SYSTEM_TOKEN     = "ajsdIgsnaloHFGis823jsdgyjTGDKijfcjk783JDUYFJyggvwejkxsnmbkjwpoj6483"
	CLIENT_ID_SECRET = "EC1FD9R5t6D3cs9CzPbgJaBJjshoVgrJrTs6U39scYzYF7HYyMlv_mal2IjLLaA9" // Auth0 RPC API client
)

func NewGoogleOautConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  SERVER_ADDRESS + "/auth/google/callback",
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
		RedirectURL:  SERVER_ADDRESS + "/auth/slack/callback",
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
		//TODO: switch to SSl
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

func NewDropboxOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		//TODO: switch to SSl
		RedirectURL: "http://localhost:8082/auth/dropbox/callback",
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
		//TODO: switch to SSl
		RedirectURL:  "http://localhost:8082/auth/box/callback",
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
		RedirectURL:  SERVER_ADDRESS + "/auth/gmail/callback",
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

// NewSystemContext System context
func NewSystemContext() context.Context {
	return metadata.NewContext(context.TODO(), map[string]string{
		"Authorization": SYSTEM_TOKEN,
	})
}

func ParseJWTToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return "", errors.InternalServerError("AuthWrapper", "Unable to retrieve metadata")
	}

	token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalServerError("Unexpected signing method", token.Header["alg"].(string))
		}

		decoded, err := base64.URLEncoding.DecodeString(CLIENT_ID_SECRET)
		if err != nil {
			return nil, err
		}

		return decoded, nil
	})

	if err != nil {
		return "", errors.Unauthorized("Token", err.Error())
	}

	if !token.Valid {
		return "", errors.Unauthorized("", "Invalid token")
	}

	return token.Claims.(jwt.MapClaims)["sub"].(string), nil
}

// NewUUID generates a random UUID according to RFC 4122
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

// Remove records (Files) from db that not longer belong to a datasource
// Compares LastSeen with the time the crawler started
// so all records with a LastSeen before will be removed from index
// file does not exists any more on datasource
func ClearIndex(e *datasource_proto.Endpoint) error {
	c := db_proto.NewDBClient(DB_SERVICE_NAME, nil)
	_, err := c.DeleteByQuery(NewSystemContext(), &db_proto.DeleteByQueryRequest{
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
