package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	_ "github.com/kazoup/platform/lib/db/config/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/errors"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"testing"
)

func TestNewDataSourceEngine(t *testing.T) {
	testData := []struct {
		url string
		err error
	}{
		{
			url: "googledrive://test",
			err: nil,
		},
		{
			url: "gmail://test",
			err: nil,
		},
		{
			url: "onedrive://test",
			err: nil,
		},
		{
			url: "slack://test",
			err: nil,
		},
		{
			url: "dropbox://test",
			err: nil,
		},
		{
			url: "box://test",
			err: nil,
		},
		{
			url: "mock",
			err: nil,
		},
		{
			url: "invalid",
			err: errors.ErrInvalidDatasourceEngine,
		},
	}

	for _, tt := range testData {
		_, err := NewDataSourceEngine(&proto_datasource.Endpoint{
			Url: tt.url,
		})

		if tt.err != err {
			t.Errorf("Expected %v, got %v", tt.err, err)
		}
	}
}

func TestGenerateEndpoint(t *testing.T) {
	endpoint, err := GenerateEndpoint(context.TODO(), proto_datasource.Endpoint{
		Url: globals.Mock,
	})

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if len(endpoint.Id) == 0 {
		t.Error("Id not found")
	}

	if len(endpoint.Index) == 0 {
		t.Error("Index not found")
	}
}

func TestReadDataSource(t *testing.T) {
	endpoint, err := ReadDataSource(context.TODO(), "test_id")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if globals.Mock != endpoint.Url {
		t.Error("Expected %v, got %v", globals.Mock, endpoint.Url)
	}
}

func TestSearchDataSources(t *testing.T) {
	expectedResult := `[{"url":"mock"}]`

	rsp, err := SearchDataSources(context.TODO(), &proto_datasource.SearchRequest{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expectedResult != rsp.Result {
		t.Errorf("Expected %v, got %v", expectedResult, rsp.Result)
	}
}

func TestSaveDataSource(t *testing.T) {
	result := SaveDataSource(context.TODO(), proto_datasource.Endpoint{
		Url: globals.Mock,
	}, "test_id")

	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestDeleteDataSource(t *testing.T) {
	result := DeleteDataSource(context.TODO(), &proto_datasource.Endpoint{
		Url: globals.Mock,
	})
	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestCreateIndexWithAlias(t *testing.T) {
	result := CreateIndexWithAlias(context.TODO(), &proto_datasource.Endpoint{
		Url: globals.Mock,
	})
	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}
