package subscriber

import (
	"encoding/json"
	"github.com/kazoup/platform/entities/srv/proto/entities"
	kazoup_context "github.com/kazoup/platform/lib/context"
	platform_errors "github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	announce "github.com/kazoup/platform/lib/protomsg/announce"
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
	sentimentHandler = new(SentimentHandler)
	srv              = micro.NewService(
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

func TestOnAudioEnrich(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Id:     "test",
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
		err := sentimentHandler.OnAudioEnrich(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnDocEnrich(t *testing.T) {
	e := &enrich.EnrichMessage{
		UserId: TEST_USER_ID,
		Id:     "test",
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
		err := sentimentHandler.OnDocEnrich(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}

func TestOnAnalyzeFile(t *testing.T) {
	e := &proto_entities.ExtractFileRequest{}

	b, err := json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var onAnalyzeFileTestData = []struct {
		ctx context.Context
		msg *announce.AnnounceMessage
		out error
	}{
		// Success
		{
			micro.NewContext(ctx, srv),
			&announce.AnnounceMessage{
				Handler: globals.HANDLER_SENTIMENT_ENRICH_FILE,
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
				Handler: globals.HANDLER_SENTIMENT_ENRICH_FILE,
				Data:    string(b),
			},
			platform_errors.ErrInvalidCtx,
		},
	}

	for _, tt := range onAnalyzeFileTestData {
		err := sentimentHandler.OnAnalyzeFile(tt.ctx, tt.msg)
		if err != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, err)
		}
	}
}
