package handler

import (
	"github.com/kazoup/platform/entities/srv/proto/entities"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/quota/mock"
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

func TestExtractFile(t *testing.T) {
	var enrichFileTestData = []struct {
		ctx         context.Context
		req         *proto_entities.ExtractFileRequest
		expectedRsp *proto_entities.ExtractFileResponse
		rsp         *proto_entities.ExtractFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_entities.ExtractFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_entities.ExtractFileResponse{
				Info: QUOTA_EXCEEDED_MSG,
			},
			&proto_entities.ExtractFileResponse{},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_entities.ExtractFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_entities.ExtractFileResponse{
				Info: "",
			},
			&proto_entities.ExtractFileResponse{},
		},
	}

	for _, tt := range enrichFileTestData {
		if err := srv.ExtractFile(tt.ctx, tt.req, tt.rsp); err != nil {
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
		req         *proto_entities.HealthRequest
		expectedRsp *proto_entities.HealthResponse
		rsp         *proto_entities.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_entities.HealthRequest{},
			&proto_entities.HealthResponse{
				Status: 200,
			},
			&proto_entities.HealthResponse{},
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
