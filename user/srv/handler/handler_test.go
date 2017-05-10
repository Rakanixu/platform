package handler

import (
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/user/srv/proto/user"
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

func TestRead(t *testing.T) {
	var readTestData = []struct {
		ctx         context.Context
		req         *proto_user.ReadRequest
		expectedRsp *proto_user.ReadResponse
		rsp         *proto_user.ReadResponse
	}{
		{
			metadata.NewContext(ctx, map[string]string{}),
			&proto_user.ReadRequest{
				Index: "test",
				Id:    "test",
			},
			&proto_user.ReadResponse{},
			&proto_user.ReadResponse{},
		},
	}

	for _, tt := range readTestData {
		if err := srv.Read(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Result != tt.rsp.Result {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp, tt.rsp)
		}
	}
}

func TestHealth(t *testing.T) {
	var healthTestData = []struct {
		ctx         context.Context
		req         *proto_user.HealthRequest
		expectedRsp *proto_user.HealthResponse
		rsp         *proto_user.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_user.HealthRequest{},
			&proto_user.HealthResponse{
				Status: 200,
			},
			&proto_user.HealthResponse{},
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
