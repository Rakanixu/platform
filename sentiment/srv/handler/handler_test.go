package handler

import (
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/quota/mock"
	"github.com/kazoup/platform/sentiment/srv/proto/sentiment"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	srv = new(Service)
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestAnalyzeFile(t *testing.T) {
	var enrinchFilesTestData = []struct {
		ctx         context.Context
		req         *proto_sentiment.AnalyzeFileRequest
		expectedRsp *proto_sentiment.AnalyzeFileResponse
		rsp         *proto_sentiment.AnalyzeFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_sentiment.AnalyzeFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_sentiment.AnalyzeFileResponse{
				Info: QUOTA_EXCEEDED_MSG,
			},
			&proto_sentiment.AnalyzeFileResponse{},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_sentiment.AnalyzeFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_sentiment.AnalyzeFileResponse{},
			&proto_sentiment.AnalyzeFileResponse{},
		},
	}

	for _, tt := range enrinchFilesTestData {
		if err := srv.AnalyzeFile(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Info != tt.rsp.Info {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Info, tt.rsp.Info)
		}
	}
}

func TestHealth(t *testing.T) {
	var healthTestData = []struct {
		ctx         context.Context
		req         *proto_sentiment.HealthRequest
		expectedRsp *proto_sentiment.HealthResponse
		rsp         *proto_sentiment.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_sentiment.HealthRequest{},
			&proto_sentiment.HealthResponse{
				Status: 200,
			},
			&proto_sentiment.HealthResponse{},
		},
	}

	for _, tt := range healthTestData {
		if err := srv.Health(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Status != tt.rsp.Status {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Status, tt.rsp.Status)
		}
	}
}
