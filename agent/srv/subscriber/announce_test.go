package subscriber

import (
	"encoding/json"
	kazoup_context "github.com/kazoup/platform/lib/context"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/file"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro"
	broker_mock "github.com/micro/go-micro/broker/mock"
	registry_mock "github.com/micro/go-micro/registry/mock"
	"golang.org/x/net/context"
	"testing"
	"time"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	announceHandler = new(AnnounceHandler)
	srv             = micro.NewService(
		micro.Name("test-service"),
		micro.Broker(broker_mock.NewBroker()),
		micro.Registry(registry_mock.NewRegistry()),
	)
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

// Announce subscriber OnSave method unit test
func TestOnSave(t *testing.T) {
	// Parse string to Time
	m, err := time.Parse(time.RFC3339, "2017-03-23T10:43:34Z")
	if err != nil {
		t.Fatal(err)
	}

	// Testing KazoupFile
	f := file.KazoupFile{
		ID:                  "a716f1408cfff7afd943acd45dcfa0a4",
		OriginalID:          "id:lXWZMx78s2AAAAAAAAAAPw",
		OriginalDownloadRef: "",
		PreviewUrl:          "/gopher.jpg",
		UserId:              "google-apps|pablo.aguirre@kazoup.com",
		Name:                "gopher.jpg",
		URL:                 "https://www.dropbox.com/home?preview=gopher.jpg",
		Modified:            m,
		FileSize:            75679,
		IsDir:               false,
		Category:            "Pictures",
		MimeType:            "",
		Depth:               0,
		FileType:            "dropbox",
		LastSeen:            1496909821,
		Access:              "private",
		DatasourceId:        "e80d54ad29d18cb62cf9bb2bb54fcfd5",
	}

	// Marshal file structure to JSON
	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	var OnSaveTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_AGENT_SAVE,
				Data:    string(b),
			},
			nil,
		},
		// Ignore msg due to topic
		{
			ctx,
			&announce.AnnounceMessage{
				Handler: "ignore-me",
			},
			nil,
		},
		//Invalid context
		{
			ctx,
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_AGENT_SAVE,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range OnSaveTestData {
		err := announceHandler.OnSave(tt.ctx, tt.msg)
		if err != tt.out {
			t.Fatal(err)
		}
	}
}
