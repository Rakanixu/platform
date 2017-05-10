package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	"github.com/kazoup/platform/document/srv/proto/document"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/custom/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/kazoup/platform/lib/wrappers"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	announceHandler = new(AnnounceHandler)
	srv             = wrappers.NewKazoupService("test-service")
	ctx             = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestOnCrawlerFinished(t *testing.T) {
	e := &proto_datasource.Endpoint{
		UserId: TEST_USER_ID,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var OnCrawlerFinishedTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.DiscoverTopic,
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
				Handler: globals.DiscoverTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range OnCrawlerFinishedTestData {
		err := announceHandler.OnCrawlerFinished(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnEnrichFile(t *testing.T) {
	m := &proto_document.EnrichFileRequest{}

	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var OnCrawlerFinishedTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_DOCUMENT_ENRICH_FILE,
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
				Handler: globals.HANDLER_DOCUMENT_ENRICH_FILE,
				Data:    "aasasdf",
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range OnCrawlerFinishedTestData {
		err := announceHandler.OnEnrichFile(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnEnrichDatasource(t *testing.T) {
	e := &proto_document.EnrichDatasourceRequest{}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var OnCrawlerFinishedTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_AUDIO_ENRICH_DATASOURCE,
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
				Handler: globals.HANDLER_DOCUMENT_ENRICH_DATASOURCE,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range OnCrawlerFinishedTestData {
		err := announceHandler.OnEnrichDatasource(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}
