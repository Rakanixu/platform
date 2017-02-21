package handler

import (
	"bytes"
	"fmt"
	proto_datasource "github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/lib/globals"
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	"github.com/kazoup/platform/lib/wrappers"
	notification_proto "github.com/kazoup/platform/notification/srv/proto/notification"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

var gcs *gcslib.GoogleCloudStorage

func init() {
	gcslib.Register()
	gcs = gcslib.NewGoogleCloudStorage()
}

//SaveDatasource call datasource-srv and save new data source
func SaveDatasource(ctx context.Context, user string, url string, token *oauth2.Token) error {
	client := wrappers.NewKazoupClient()

	srvReq := client.NewRequest(
		globals.DATASOURCE_SERVICE_NAME,
		"DataSource.Create",
		&proto_datasource.CreateRequest{
			Endpoint: &proto_datasource.Endpoint{
				UserId:          user,
				Url:             url,
				LastScan:        0,
				LastScanStarted: 0,
				CrawlerRunning:  false,
				Token: &proto_datasource.Token{
					AccessToken:  token.AccessToken,
					TokenType:    token.TokenType,
					RefreshToken: token.RefreshToken,
					Expiry:       token.Expiry.Unix(),
				},
			},
		},
	)

	srvRes := &proto_datasource.CreateResponse{}

	if err := client.Call(ctx, srvReq, srvRes); err != nil {
		return errors.New("auth.web.helper.DataSource.Create", err.Error(), 500)
	}

	return nil
}

// SaveTmpToken saves JWT in GCS
func SaveTmpToken(uuid, jwt string) error {
	_, err := globals.ParseJWTToken(jwt) // Parse JWT to be sure was signed by us
	if err != nil {
		return err
	}

	// We save uuid - jwt token pair in GCS
	if err := gcs.Upload(bytes.NewBufferString(jwt), globals.TMP_TOKEN_BUCKET, uuid); err != nil {
		return err
	}

	return nil
}

// RetrieveUserAndContextFromUUID retrieves userId and context from GCS
func RetrieveUserAndContextFromUUID(uuid string) (string, context.Context, error) {
	// Retrieve JWT associated with uuid
	rd, err := gcs.Download(globals.TMP_TOKEN_BUCKET, string(uuid))
	if err != nil {
		return "", nil, err
	}
	defer rd.Close()

	jwt, err := ioutil.ReadAll(rd)
	if err != nil {
		return "", nil, err
	}

	// Parse JWT as a way to validate it, and retrieve user_id associated with that JWT
	uID, err := globals.ParseJWTToken(string(jwt))
	if err != nil {
		return "", nil, err
	}

	return uID, globals.NewContextFromJWT(string(jwt)), nil
}

//PublishNotification send data source created notification
func PublishNotification(uID string) error {
	n := &notification_proto.NotificationMessage{
		Info:   "Datasource created succesfully",
		Method: globals.NOTIFY_REFRESH_DATASOURCES,
		UserId: string(uID),
	}

	// Publish scan topic, crawlers should pick up message and start scanning
	return client.DefaultClient.Publish(
		globals.NewSystemContext(),
		client.DefaultClient.NewPublication(globals.NotificationTopic, n),
	)
}

// CloseBrowserWindow loads app in settings page, and the close that window
func CloseBrowserWindow(w http.ResponseWriter, r *http.Request) {
	// Close window
	fmt.Fprintf(w, "%s", `
		<script>
		'use stric';
			(function() {
				window.close();
			}());
		</script>
	`)
}
