package globals

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	datasource_proto "github.com/kazoup/platform/datasource/srv/proto/datasource"
	db_proto "github.com/kazoup/platform/db/srv/proto/db"
	kazoup_context "github.com/kazoup/platform/lib/context"
	"github.com/micro/go-micro/client"
	micro_errors "github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/slack"
	"io"
	"log"
	"time"
)

const (
	NAMESPACE                      string = "com.kazoup"
	DB_SERVICE_NAME                string = NAMESPACE + ".srv.db"
	SEARCH_SERVICE_NAME            string = NAMESPACE + ".srv.search"
	DATASOURCE_SERVICE_NAME        string = NAMESPACE + ".srv.datasource"
	CRAWLER_SERVICE_NAME           string = NAMESPACE + ".srv.crawler"
	NOTIFICATION_SERVICE_NAME      string = NAMESPACE + ".srv.notification"
	FILE_SERVICE_NAME              string = NAMESPACE + ".srv.file"
	QUOTA_SERVICE_NAME             string = NAMESPACE + ".srv.quota"
	PROFILE_SERVICE_NAME           string = NAMESPACE + ".srv.profile"
	MONITOR_SERVICE_NAME           string = NAMESPACE + ".srv.monitor"
	THUMBNAIL_SERVICE_NAME         string = NAMESPACE + ".srv.thumbnail"
	AUDIOENRICH_SERVICE_NAME       string = NAMESPACE + ".srv.audioenrich"
	DOCENRICH_SERVICE_NAME         string = NAMESPACE + ".srv.docenrich"
	IMGENRICH_SERVICE_NAME         string = NAMESPACE + ".srv.imgenrich"
	TEXTANALYZER_SERVICE_NAME      string = NAMESPACE + ".srv.textanalyzer"
	SENTIMENTANALYZER_SERVICE_NAME string = NAMESPACE + ".srv.sentimentanalyzer"

	HANDLER_DATASOURCE_CREATE  = DATASOURCE_SERVICE_NAME + ".DataSource.create"
	HANDLER_DATASOURCE_DELETE  = DATASOURCE_SERVICE_NAME + ".DataSource.Delete"
	HANDLER_DATASOURCE_SEARCH  = DATASOURCE_SERVICE_NAME + ".DataSource.Search"
	HANDLER_DATASOURCE_SCAN    = DATASOURCE_SERVICE_NAME + ".DataSource.Scan"
	HANDLER_DATASOURCE_SCANALL = DATASOURCE_SERVICE_NAME + ".DataSource.ScanAll"

	AnnounceTopic           string = NAMESPACE + ".topic.announce"
	FilesTopic              string = NAMESPACE + ".topic.files"
	SlackChannelsTopic      string = NAMESPACE + ".topic.slackchannels"
	SlackUsersTopic         string = NAMESPACE + ".topic.slackusers"
	ScanTopic               string = NAMESPACE + ".topic.scan"
	DocEnrichTopic          string = NAMESPACE + ".topic.docenrich"
	ImgEnrichTopic          string = NAMESPACE + ".topic.imgenrich"
	ThumbnailTopic          string = NAMESPACE + ".topic.thumbnail"
	AudioEnrichTopic        string = NAMESPACE + ".topic.audioenrich"
	SentimentEnrichTopic    string = NAMESPACE + ".topic.sentiment"
	ExtractEntitiesTopic    string = NAMESPACE + ".topic.extractentities"
	CrawlerStartedTopic     string = NAMESPACE + ".topic.crawlerstarted"
	CrawlerFinishedTopic    string = NAMESPACE + ".topic.crawlerfinished"
	NotificationTopic       string = NAMESPACE + ".topic.notification"
	NotificationProxyTopic  string = NAMESPACE + ".topic.notificationproxy"
	DeleteBucketTopic       string = NAMESPACE + ".topic.deletebucket"
	DeleteFileInBucketTopic string = NAMESPACE + ".topic.deletefileinbucket"

	IndexDatasources  = "datasources"
	IndexHelper       = "files_helper"
	FileType          = "file"
	FoldeType         = "folder"
	UserType          = "user"
	ChannelType       = "channel"
	FileTypeFile      = "files"
	FileTypeDirectory = "directories"
	TypeDatasource    = "datasource"

	Local       = "local"
	Slack       = "slack"
	GoogleDrive = "googledrive"
	Gmail       = "gmail"
	OneDrive    = "onedrive"
	Dropbox     = "dropbox"
	Box         = "box"

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

	DISCOVERY_DELAY_MS  = 10 * time.Millisecond
	PUBLISHING_DELAY_MS = 20 * time.Millisecond
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

// DBAccess flags access to DB
func DBAccess(ctx context.Context) error {
	md, _ := metadata.FromContext(ctx)

	if len(md["X-Kazoup-Token"]) == 0 {
		return micro_errors.Forbidden("com.kazoup.srv.db", "No scope")
	}

	if md["X-Kazoup-Token"] != DB_ACCESS_TOKEN {
		return micro_errors.Forbidden("com.kazoup.srv.db", "Invalid scope")
	}

	return nil
}

// ParseUserIdFromContext returns user_id from context
func ParseUserIdFromContext(ctx context.Context) (string, error) {
	if ctx.Value(kazoup_context.UserIdCtxKey{}) == nil {
		return "", micro_errors.Unauthorized("ParseUserIdFromContext", "Unable to retrieve user from context")
	}

	id := string(ctx.Value(kazoup_context.UserIdCtxKey{}).(kazoup_context.UserIdCtxValue))
	if len(id) == 0 {
		return "", micro_errors.Unauthorized("ParseUserIdFromContext", "No user for given context")
	}

	return id, nil
}

func ParseRolesFromContext(ctx context.Context) ([]string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return []string{}, errors.New("Unable to retrieve metadata")
	}

	if len(md["Authorization"]) == 0 {
		return []string{}, errors.New("No Auth header")
	}

	// We will read claim to know if public user, or paying or whatever
	token, err := jwt.Parse(md["Authorization"], func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		decoded, err := base64.URLEncoding.DecodeString(CLIENT_ID_SECRET)
		if err != nil {
			return nil, err
		}

		return decoded, nil
	})
	if err != nil {
		return []string{}, err
	}

	if token.Claims.(jwt.MapClaims)["roles"] == nil {
		return []string{}, errors.New("Roles not found.")
	}

	var roles []string
	for _, v := range token.Claims.(jwt.MapClaims)["roles"].([]interface{}) {
		roles = append(roles, v.(string))
	}

	return roles, nil
}

// ParseJWTToken validates JWT and returns user_id claim
func ParseJWTToken(str string) (string, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, micro_errors.InternalServerError("Unexpected signing method", token.Header["alg"].(string))
		}

		decoded, err := base64.URLEncoding.DecodeString(CLIENT_ID_SECRET)
		if err != nil {
			return nil, err
		}

		return decoded, nil
	})

	if err != nil {
		return "", micro_errors.Unauthorized("Token", err.Error())
	}

	if !token.Valid {
		return "", micro_errors.Unauthorized("", "Invalid token")
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
// Also deletes thumbs that does not exists any more on index
func ClearIndex(ctx context.Context, c client.Client, e *datasource_proto.Endpoint) error {
	// Clean the index after all messages have been published
	delReq := &db_proto.DeleteByQueryRequest{
		Indexes:  []string{e.Index},
		Types:    []string{"file"},
		LastSeen: e.LastScanStarted,
	}
	srvReq := c.NewRequest(
		DB_SERVICE_NAME,
		"DB.DeleteByQuery",
		delReq,
	)
	srvRes := &db_proto.DeleteByQueryResponse{}
	if err := c.Call(ctx, srvReq, srvRes); err != nil {
		log.Printf("Error globals.ClearIndex -  %s", err)
		return err
	}

	return nil
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

var fileTypeDict = struct {
	m map[string]map[string]string
}{
	m: map[string]map[string]string{
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

func GetMimeType(fileSystemType, fileType string) string {
	// Be sure to not panic if input not in map
	if fileTypeDict.m[fileSystemType] != nil {
		if len(fileTypeDict.m[fileSystemType][fileType]) > 0 {
			return fileTypeDict.m[fileSystemType][fileType]
		}
	}

	return ""
}

func GoogleDriveExportAs(originalMimeType string) string {
	//https://developers.google.com/drive/v3/web/integrate-open#open_and_convert_google_docs_in_your_app
	switch originalMimeType {
	case GOOGLE_DRIVE_DOCUMENT:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case GOOGLE_DRIVE_PRESETATION:
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case GOOGLE_DRIVE_SPREADSHEET:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}

	return ""
}

func GetDocumentTemplate(fileType string, fullName bool) string {
	var tmp string

	switch fileType {
	case DOCUMENT:
		tmp = "docx"
	case PRESENTATION:
		tmp = "pptx"
	case SPREADSHEET:
		tmp = "xlsx"
	case TEXT:
		tmp = "txt"
	}

	if fullName {
		tmp = fmt.Sprintf("%s.%s", tmp, tmp)
	}

	return tmp
}

// Encrypt slice of bytes
func Encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

// Decrypt slice of bytes
func Decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
