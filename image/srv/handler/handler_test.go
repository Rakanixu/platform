package handler

import (
	"github.com/kazoup/platform/image/srv/proto/image"
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
		req         *proto_image.EnrichFileRequest
		expectedRsp *proto_image.EnrichFileResponse
		rsp         *proto_image.EnrichFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_image.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_image.EnrichFileResponse{
				Info: QUOTA_EXCEEDED_MSG,
			},
			&proto_image.EnrichFileResponse{},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_image.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_image.EnrichFileResponse{
				Info: "",
			},
			&proto_image.EnrichFileResponse{},
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
		req         *proto_image.EnrichDatasourceRequest
		expectedRsp *proto_image.EnrichDatasourceResponse
		rsp         *proto_image.EnrichDatasourceResponse
	}{
		{
			context.TODO(),
			&proto_image.EnrichDatasourceRequest{
				Id: "test_id",
			},
			&proto_image.EnrichDatasourceResponse{},
			&proto_image.EnrichDatasourceResponse{},
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
	var healthTestData = []struct {
		ctx         context.Context
		req         *proto_image.HealthRequest
		expectedRsp *proto_image.HealthResponse
		rsp         *proto_image.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_image.HealthRequest{},
			&proto_image.HealthResponse{
				Status: 200,
			},
			&proto_image.HealthResponse{},
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
