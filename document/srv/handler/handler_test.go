package handler

import (
	"github.com/kazoup/platform/document/srv/proto/document"
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
	var enrichFileTestData = []struct {
		ctx         context.Context
		req         *proto_document.EnrichFileRequest
		expectedRsp *proto_document.EnrichFileResponse
		rsp         *proto_document.EnrichFileResponse
	}{
		// Quota has been excedded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "true",
			}),
			&proto_document.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_document.EnrichFileResponse{
				Info: QUOTA_EXCEEDED_MSG,
			},
			&proto_document.EnrichFileResponse{},
		},
		// Quota has not been exceeded
		{
			metadata.NewContext(ctx, map[string]string{
				"Quota-Exceeded": "false",
			}),
			&proto_document.EnrichFileRequest{
				Index: "test_index",
				Id:    "test_id",
			},
			&proto_document.EnrichFileResponse{
				Info: "",
			},
			&proto_document.EnrichFileResponse{},
		},
	}

	for _, tt := range enrichFileTestData {
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
		req         *proto_document.EnrichDatasourceRequest
		expectedRsp *proto_document.EnrichDatasourceResponse
		rsp         *proto_document.EnrichDatasourceResponse
	}{
		{
			context.TODO(),
			&proto_document.EnrichDatasourceRequest{
				Id: "test_id",
			},
			&proto_document.EnrichDatasourceResponse{},
			&proto_document.EnrichDatasourceResponse{},
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
		req         *proto_document.HealthRequest
		expectedRsp *proto_document.HealthResponse
		rsp         *proto_document.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_document.HealthRequest{},
			&proto_document.HealthResponse{
				Status: 200,
			},
			&proto_document.HealthResponse{},
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
