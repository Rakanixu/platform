package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/audio/srv/proto/audio"
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/custom/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	"github.com/micro/go-micro"
	broker_mock "github.com/micro/go-micro/broker/mock"
	"github.com/micro/go-micro/metadata"
	registry_mock "github.com/micro/go-micro/registry/mock"
	"golang.org/x/net/context"
	"testing"
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
	m := &proto_audio.EnrichFileRequest{}

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
				Handler: globals.HANDLER_IMAGE_ENRICH_FILE,
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
				Handler: globals.HANDLER_IMAGE_ENRICH_FILE,
				Data:    string(b),
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
	e := &proto_audio.EnrichDatasourceRequest{}

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
			micro.NewContext(metadata.NewContext(ctx, map[string]string{
				"Wanted-Type": globals.TypeDatasource,
			}), srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_IMAGE_ENRICH_DATASOURCE,
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
				Handler: globals.HANDLER_IMAGE_ENRICH_DATASOURCE,
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
