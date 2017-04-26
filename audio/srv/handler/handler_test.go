package handler

import (
	"github.com/kazoup/platform/audio/srv/proto/audio"
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

func TestEnrichFile(t *testing.T) {
	var enrinchFilesTestData = []struct {
		ctx         context.Context
		req         *proto_audio.EnrichFileRequest
		expectedRsp *proto_audio.EnrichFileResponse
		rsp         *proto_audio.EnrichFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_audio.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_audio.EnrichFileResponse{
				Info: QUOTA_EXCEEDED_MSG,
			},
			&proto_audio.EnrichFileResponse{},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_audio.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_audio.EnrichFileResponse{
				Info: "",
			},
			&proto_audio.EnrichFileResponse{},
		},
	}

	for _, tt := range enrinchFilesTestData {
		if err := srv.EnrichFile(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Info != tt.rsp.Info {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Info, tt.rsp.Info)
		}
	}
}

func TestEnrichDatasource(t *testing.T) {
	var enrinchDatasourceTestData = []struct {
		ctx         context.Context
		req         *proto_audio.EnrichDatasourceRequest
		expectedRsp *proto_audio.EnrichDatasourceResponse
		rsp         *proto_audio.EnrichDatasourceResponse
	}{
		{
			context.TODO(),
			&proto_audio.EnrichDatasourceRequest{
				Id: "test_id",
			},
			&proto_audio.EnrichDatasourceResponse{},
			&proto_audio.EnrichDatasourceResponse{},
		},
	}

	for _, tt := range enrinchDatasourceTestData {
		if err := srv.EnrichDatasource(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Info != tt.rsp.Info {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Info, tt.rsp.Info)
		}
	}
}

func TestHealth(t *testing.T) {
	var ehealthTestData = []struct {
		ctx         context.Context
		req         *proto_audio.HealthRequest
		expectedRsp *proto_audio.HealthResponse
		rsp         *proto_audio.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_audio.HealthRequest{},
			&proto_audio.HealthResponse{
				Status: 200,
			},
			&proto_audio.HealthResponse{},
		},
	}

	for _, tt := range ehealthTestData {
		if err := srv.Health(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Status != tt.rsp.Status {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Status, tt.rsp.Status)
		}
	}
}
