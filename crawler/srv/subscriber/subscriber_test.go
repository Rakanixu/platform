package subscriber

import (
	"github.com/kazoup/platform/datasource/srv/proto/datasource"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	"github.com/kazoup/platform/lib/globals"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"testing"
)

var (
	taskHandler = new(TaskHandler)
)

func TestScans(t *testing.T) {
	var scansTestData = []struct {
		ctx      context.Context
		endpoint *proto_datasource.Endpoint
		result   error
	}{
		{
			micro.NewContext(ctx, srv),
			&proto_datasource.Endpoint{
				Url: globals.Mock,
			},
			nil,
		},
	}

	for _, tt := range scansTestData {
		result := taskHandler.Scans(tt.ctx, tt.endpoint)

		if tt.result != result {
			t.Error("Expected %v, got %v", tt.result, result)
		}
	}
}
