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
	workers := 3
	th := &taskHandler{
		enrichMsgChan: make(chan enrichMsgChan, 1000000),
		workers:       workers,
	}

	for i := 0; i < th.workers; i++ {
		go th.queueListener(i)
	}

	var queueListenerTestData = []struct {
		msg enrichMsgChan
	}{
		{
			enrichMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
		{
			enrichMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
		{
			enrichMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
	}

	for _, tt := range queueListenerTestData {
		th.enrichMsgChan <- tt.msg
	}

	if len(queueListenerTestData) != len(th.enrichMsgChan) {
		t.Error("Expected %v, got %v", len(queueListenerTestData), len(th.enrichMsgChan))
	}
}

func TeststartWorkers(t *testing.T) {
	workers := 5
	th := &taskHandler{
		enrichMsgChan: make(chan enrichMsgChan, 1000000),
		workers:       workers,
	}

	var queueListenerTestData = []struct {
		msg enrichMsgChan
	}{
		{
			enrichMsgChan{},
		},
		{
			enrichMsgChan{},
		},
		{
			enrichMsgChan{},
		},
	}

	startWorkers(th)

	for _, tt := range queueListenerTestData {
		th.enrichMsgChan <- tt.msg
	}

	if len(queueListenerTestData) != len(th.enrichMsgChan) {
		t.Error("Expected %v, got %v", len(queueListenerTestData), len(th.enrichMsgChan))
	}
}

func TestprocessEnrichMsg(t *testing.T) {
	var enrichMsgTestData = []struct {
		msg    enrichMsgChan
		result error
	}{
		{
			enrichMsgChan{},
			nil,
		},
		/*		{
					enrichMsgChan{},
					nil,
				},
				{
					enrichMsgChan{},
					nil,
				},*/
	}

	for _, tt := range enrichMsgTestData {
		result := processEnrichMsg(tt.msg)

		if tt.result != result {
			t.Errorf("Expected '%v', got: '%v'", tt.result, result)
		}
	}
}
