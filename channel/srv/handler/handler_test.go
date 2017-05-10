package handler

import (
	"github.com/kazoup/platform/channel/srv/proto/channel"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
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
		req         *proto_channel.ReadRequest
		expectedRsp *proto_channel.ReadResponse
		rsp         *proto_channel.ReadResponse
	}{
		{
			metadata.NewContext(ctx, map[string]string{}),
			&proto_channel.ReadRequest{
				Index: "test",
				Id:    "test",
			},
			&proto_channel.ReadResponse{},
			&proto_channel.ReadResponse{},
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
		req         *proto_channel.HealthRequest
		expectedRsp *proto_channel.HealthResponse
		rsp         *proto_channel.HealthResponse
	}{
		// Assert service returns HTTP 200 OK
		{
			context.TODO(),
			&proto_channel.HealthRequest{},
			&proto_channel.HealthResponse{
				Status: 200,
			},
			&proto_channel.HealthResponse{},
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
