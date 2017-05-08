package subscriber

import (
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"testing"
)

func TestNewTaskHandler(t *testing.T) {
	workers := 1
	th := NewTaskHandler(workers)

	if th.workers != workers {
		t.Errorf("Expected %v, got %v", workers, th.workers)
	}
}

func TestTaskHandler_Enrich(t *testing.T) {
	workers := 1
	th := NewTaskHandler(workers)

	var enrichTestData = []struct {
		ctx    context.Context
		msg    *enrich_proto.EnrichMessage
		result error
	}{
		{
			metadata.NewContext(ctx, map[string]string{}),
			&enrich_proto.EnrichMessage{},
			nil,
		},
	}

	for _, tt := range enrichTestData {
		result := th.Enrich(tt.ctx, tt.msg)

		if tt.result != result {
			t.Errorf("Expected '%v', got: '%v'", tt.result, result)
		}
	}
}

func TestTaskHandler_queueListener(t *testing.T) {

}

func TeststartWorkers(t *testing.T) {

}

func TestprocessEnrichMsg(t *testing.T) {

}
