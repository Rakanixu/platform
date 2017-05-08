package handler

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/globals"
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	"golang.org/x/net/context"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	srv = new(service)
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestCreate(t *testing.T) {
	var createTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.CreateRequest
		expectedRsp *proto_datasource.CreateResponse
		rsp         *proto_datasource.CreateResponse
	}{
		{
			ctx,
			&proto_datasource.CreateRequest{
				&proto_datasource.Endpoint{
					Url: globals.Mock,
				},
			},
			&proto_datasource.CreateResponse{},
			&proto_datasource.CreateResponse{},
		},
	}

	for _, tt := range createTestData {
		if err := srv.Create(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp.Response != tt.rsp.Response {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp.Response, tt.rsp.Response)
		}
	}
}

func TestRead(t *testing.T) {
	var readTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.ReadRequest
		expectedRsp *proto_datasource.ReadResponse
		rsp         *proto_datasource.ReadResponse
	}{
		{
			ctx,
			&proto_datasource.ReadRequest{
				Id: "test",
			},
			&proto_datasource.ReadResponse{
				Result: `{"url":"mock"}`,
			},
			&proto_datasource.ReadResponse{},
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
		req         *proto_datasource.DeleteRequest
		expectedRsp *proto_datasource.DeleteResponse
		rsp         *proto_datasource.DeleteResponse
	}{
		{
			ctx,
			&proto_datasource.DeleteRequest{
				Id: "test",
			},
			&proto_datasource.DeleteResponse{},
			&proto_datasource.DeleteResponse{},
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
	var searchTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.SearchRequest
		expectedRsp *proto_datasource.SearchResponse
		rsp         *proto_datasource.SearchResponse
	}{
		{
			ctx,
			&proto_datasource.SearchRequest{
				From: 0,
				Size: 9999,
			},
			&proto_datasource.SearchResponse{
				Result: `[{"url":"mock"}]`,
			},
			&proto_datasource.SearchResponse{},
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

func TestScan(t *testing.T) {
	var scanTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.ScanRequest
		expectedRsp *proto_datasource.ScanResponse
		rsp         *proto_datasource.ScanResponse
	}{
		{
			ctx,
			&proto_datasource.ScanRequest{
				Id: "test",
			},
			&proto_datasource.ScanResponse{},
			&proto_datasource.ScanResponse{},
		},
	}

	for _, tt := range scanTestData {
		if err := srv.Scan(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp != tt.rsp {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp, tt.rsp)
		}
	}
}

func TestScanAll(t *testing.T) {
	var scanAllTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.ScanAllRequest
		expectedRsp *proto_datasource.ScanAllResponse
		rsp         *proto_datasource.ScanAllResponse
	}{
		{
			ctx,
			&proto_datasource.ScanAllRequest{
				DatasourcesId: []string{"test"},
			},
			&proto_datasource.ScanAllResponse{},
			&proto_datasource.ScanAllResponse{},
		},
	}

	for _, tt := range scanAllTestData {
		if err := srv.ScanAll(tt.ctx, tt.req, tt.rsp); err != nil {
			t.Fatal(err)
		}

		if tt.expectedRsp != tt.rsp {
			t.Errorf("Expected '%v', got: '%v'", tt.expectedRsp, tt.rsp)
		}
	}
}

func TestHealth(t *testing.T) {
	var healthTestData = []struct {
		ctx         context.Context
		req         *proto_datasource.HealthRequest
		expectedRsp *proto_datasource.HealthResponse
		rsp         *proto_datasource.HealthResponse
	}{
		{
			context.TODO(),
			&proto_datasource.HealthRequest{},
			&proto_datasource.HealthResponse{
				Status: 200,
			},
			&proto_datasource.HealthResponse{},
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
