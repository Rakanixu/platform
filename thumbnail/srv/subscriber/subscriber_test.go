package subscriber

import (
	kazoup_context "github.com/kazoup/platform/lib/context"
	_ "github.com/kazoup/platform/lib/db/operations/mock"
	enrich_proto "github.com/kazoup/platform/lib/protomsg/enrich"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/net/context"
	"testing"
)

const (
	TEST_USER_ID = "test_user"
)

var (
	ctx = context.WithValue(
		context.TODO(),
		kazoup_context.UserIdCtxKey{},
		kazoup_context.UserIdCtxValue(TEST_USER_ID),
	)
)

func TestNewTaskHandler(t *testing.T) {
	workers := 1
	th := NewTaskHandler(workers)

	if th.workers != workers {
		t.Errorf("Expected %v, got %v", workers, th.workers)
	}
}

func TestTaskHandler_Thumbnail(t *testing.T) {
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
		result := th.Thumbnail(tt.ctx, tt.msg)

		if tt.result != result {
			t.Errorf("Expected '%v', got: '%v'", tt.result, result)
		}
	}
}

func TestTaskHandler_queueListener(t *testing.T) {
	workers := 3
	th := &taskHandler{
		thumbnailMsgChan: make(chan thumbnailMsgChan, 1000000),
		workers:          workers,
	}

	for i := 0; i < th.workers; i++ {
		go th.queueListener(i)
	}

	var queueListenerTestData = []struct {
		msg thumbnailMsgChan
	}{
		{
			thumbnailMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
		{
			thumbnailMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
		{
			thumbnailMsgChan{
				msg: &enrich_proto.EnrichMessage{},
			},
		},
	}

	for _, tt := range queueListenerTestData {
		th.thumbnailMsgChan <- tt.msg
	}

	if len(queueListenerTestData) != len(th.thumbnailMsgChan) {
		t.Error("Expected %v, got %v", len(queueListenerTestData), len(th.thumbnailMsgChan))
	}
}

func TeststartWorkers(t *testing.T) {
	workers := 5
	th := &taskHandler{
		thumbnailMsgChan: make(chan thumbnailMsgChan, 1000000),
		workers:          workers,
	}

	var queueListenerTestData = []struct {
		msg thumbnailMsgChan
	}{
		{
			thumbnailMsgChan{},
		},
		{
			thumbnailMsgChan{},
		},
		{
			thumbnailMsgChan{},
		},
	}

	startWorkers(th)

	for _, tt := range queueListenerTestData {
		th.thumbnailMsgChan <- tt.msg
	}

	if len(queueListenerTestData) != len(th.thumbnailMsgChan) {
		t.Error("Expected %v, got %v", len(queueListenerTestData), len(th.thumbnailMsgChan))
	}
}

func TestprocessEnrichMsg(t *testing.T) {
	var enrichMsgTestData = []struct {
		msg    thumbnailMsgChan
		result error
	}{
		{
			thumbnailMsgChan{},
			nil,
		},
	}

	for _, tt := range enrichMsgTestData {
		result := processThumbnailMsg(tt.msg)

		if tt.result != result {
			t.Errorf("Expected '%v', got: '%v'", tt.result, result)
		}
	}
}
