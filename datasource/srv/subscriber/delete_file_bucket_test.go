package subscriber

import (
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	deletefilebucket "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	"golang.org/x/net/context"
	"testing"
)

var (
	deleteFileInBuckethandler = NewDeleteFileInBucketHandler()
)

func TestSubscribeDeleteFileInBucket(t *testing.T) {
	var subscribeDeleteFileInBucketTestData = []struct {
		ctx context.Context
		msg *deletefilebucket.DeleteFileInBucketMsg
		out error
	}{
		{
			ctx,
			&deletefilebucket.DeleteFileInBucketMsg{
				Index:  "test",
				FileId: "test",
			},
			nil,
		},
	}

	for _, tt := range subscribeDeleteFileInBucketTestData {
		result := deleteFileInBuckethandler.SubscribeDeleteFileInBucket(tt.ctx, tt.msg)
		if result != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, result)
		}
	}
}
