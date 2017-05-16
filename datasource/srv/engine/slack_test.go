package engine

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	_ "github.com/kazoup/platform/lib/db/config/mock"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/globals"
	"golang.org/x/net/context"
	"testing"
)

func TestSlack_Validate(t *testing.T) {
	b := Slack{
		Endpoint: proto_datasource.Endpoint{
			Url: globals.Mock,
		},
	}

	endpoint, err := b.Validate(context.TODO(), "[]")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(endpoint.Id) == 0 {
		t.Error("Id not found")
	}

	if len(endpoint.Index) == 0 {
		t.Error("Index not found")
	}
}

func TestSlack_Save(t *testing.T) {
	b := Slack{
		Endpoint: proto_datasource.Endpoint{
			Url: globals.Mock,
		},
	}

	result := b.Save(context.TODO(), proto_datasource.Endpoint{
		Url: globals.Mock,
	}, "test_id")

	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestSlack_Delete(t *testing.T) {
	b := Slack{
		Endpoint: proto_datasource.Endpoint{
			Url: globals.Mock,
		},
	}

	result := b.Delete(context.TODO())
	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}

func TestSlack_CreateIndexWithAlias(t *testing.T) {
	b := Slack{
		Endpoint: proto_datasource.Endpoint{
			Url: globals.Mock,
		},
	}

	result := b.CreateIndexWithAlias(context.TODO())
	if result != nil {
		t.Errorf("Unexpected error: %v", result)
	}
}
