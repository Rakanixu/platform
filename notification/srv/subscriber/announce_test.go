package subscriber

import (
	"encoding/json"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
	crawler "github.com/kazoup/platform/lib/protomsg/crawler"
	enrich "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro"
	broker_mock "github.com/micro/go-micro/broker/mock"
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

func TestOnDocEnrich(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Notify: true,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onDocEnrichTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.DocEnrichTopic,
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
				Handler: globals.DocEnrichTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onDocEnrichTestData {
		err := announceHandler.OnDocEnrich(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnImgEnrich(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Notify: true,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onImgEnrichTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.ImgEnrichTopic,
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
				Handler: globals.ImgEnrichTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onImgEnrichTestData {
		err := announceHandler.OnImgEnrich(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnAudioEnrich(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Notify: true,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onAudioEnrichTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.AudioEnrichTopic,
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
				Handler: globals.AudioEnrichTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onAudioEnrichTestData {
		err := announceHandler.OnAudioEnrich(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnSentimentExtraction(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Notify: true,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onAudioEnrichTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.SentimentEnrichTopic,
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
				Handler: globals.SentimentEnrichTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onAudioEnrichTestData {
		err := announceHandler.OnSentimentExtraction(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnEntitiesExtraction(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Notify: true,
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onAudioEnrichTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.ExtractEntitiesTopic,
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
				Handler: globals.ExtractEntitiesTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onAudioEnrichTestData {
		err := announceHandler.OnEntitiesExtraction(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnCrawlerFinished(t *testing.T) {
	e := &crawler.CrawlerFinishedMessage{
		DatasourceId: "test",
	}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onCrawlerFinishedTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.DiscoveryFinishedTopic,
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
				Handler: globals.DiscoveryFinishedTopic,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onCrawlerFinishedTestData {
		result := announceHandler.OnCrawlerFinished(tt.ctx, tt.msg)
		if result != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, result)
		}
	}
}

func TestOnFileDeleted(t *testing.T) {
	var onFileDeletedTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_FILE_DELETE,
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
				Handler: globals.HANDLER_FILE_DELETE,
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onFileDeletedTestData {
		result := announceHandler.OnFileDeleted(tt.ctx, tt.msg)
		if result != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, result)
		}
	}
}
