package handler

import (
	"encoding/json"
	"github.com/kazoup/platform/file/srv/proto/file"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/bulk/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/file"
	"golang.org/x/net/context"
	"log"
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

func TestCreate(t *testing.T) {
	b, err := json.Marshal(file.NewKazoupFileFromMockFile())
	if err != nil {
		log.Fatal(err)
	}

	var createTestData = []struct {
		ctx         context.Context
		req         *proto_file.CreateRequest
		expectedRsp *proto_file.CreateResponse
		rsp         *proto_file.CreateResponse
	}{
		{
			ctx,
			&proto_file.CreateRequest{
				DatasourceId: "test",
			},
			&proto_file.CreateResponse{
				Data: string(b),
			},
			&proto_file.CreateResponse{},
		},
	}

	for _, tt := range createTestData {
		if err := srv.Create(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Data != tt.rsp.Data {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Data, tt.rsp.Data)
		}
	}
}

func TestRead(t *testing.T) {
	b, err := json.Marshal(file.NewKazoupFileFromMockFile())
	if err != nil {
		log.Fatal(err)
	}

	var readTestData = []struct {
		ctx         context.Context
		req         *proto_file.ReadRequest
		expectedRsp *proto_file.ReadResponse
		rsp         *proto_file.ReadResponse
	}{
		{
			ctx,
			&proto_file.ReadRequest{
				Index: "test",
				Id:    "test",
			},
			&proto_file.ReadResponse{
				Result: string(b),
			},
			&proto_file.ReadResponse{},
		},
	}

	for _, tt := range readTestData {
		if err := srv.Read(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Result != tt.rsp.Result {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Result, tt.rsp.Result)
		}
	}
}

func TestDelete(t *testing.T) {
	var deleteTestData = []struct {
		ctx         context.Context
		req         *proto_file.DeleteRequest
		expectedRsp *proto_file.DeleteResponse
		rsp         *proto_file.DeleteResponse
	}{
		{
			ctx,
			&proto_file.DeleteRequest{
				DatasourceId: "test",
			},
			&proto_file.DeleteResponse{},
			&proto_file.DeleteResponse{},
		},
	}

	for _, tt := range deleteTestData {
		if err := srv.Delete(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp != tt.rsp {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp, tt.rsp)
		}
	}
}

func TestSearch(t *testing.T) {
	b, err := json.Marshal(file.NewKazoupFileFromMockFile())
	if err != nil {
		log.Fatal(err)
	}

	var searchTestData = []struct {
		ctx         context.Context
		req         *proto_file.SearchRequest
		expectedRsp *proto_file.SearchResponse
		rsp         *proto_file.SearchResponse
	}{
		{
			ctx,
			&proto_file.SearchRequest{
				From:  0,
				Size:  9999,
				Index: "test",
			},
			&proto_file.SearchResponse{
				Result: `[` + string(b) + `]`,
			},
			&proto_file.SearchResponse{},
		},
	}

	for _, tt := range searchTestData {
		if err := srv.Search(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Result != tt.rsp.Result {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Result, tt.rsp.Result)
		}
	}
}

func TestShare(t *testing.T) {
	var scanTestData = []struct {
		ctx         context.Context
		req         *proto_file.ShareRequest
		expectedRsp *proto_file.ShareResponse
		rsp         *proto_file.ShareResponse
	}{
		{
			ctx,
			&proto_file.ShareRequest{
				OriginalId:   "test",
				DatasourceId: "test",
			},
			&proto_file.ShareResponse{},
			&proto_file.ShareResponse{},
		},
	}

	for _, tt := range scanTestData {
		if err := srv.Share(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.PublicUrl != tt.rsp.PublicUrl {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.PublicUrl, tt.rsp.PublicUrl)
		}
	}
}

func TestHealth(t *testing.T) {
	var healthTestData = []struct {
		ctx         context.Context
		req         *proto_file.HealthRequest
		expectedRsp *proto_file.HealthResponse
		rsp         *proto_file.HealthResponse
	}{
		{
			context.TODO(),
			&proto_file.HealthRequest{},
			&proto_file.HealthResponse{
				Status: 200,
			},
			&proto_file.HealthResponse{},
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
