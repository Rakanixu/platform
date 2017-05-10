package subscriber

import (
	_ "github.com/kazoup/platform/lib/objectstorage/mock"
	deletebucket "github.com/kazoup/platform/lib/protomsg/deletebucket"
	"golang.org/x/net/context"
	"testing"
)

var (
	deleteBuckethandler = NewDeleteBucketHandler()
)

func TestSubscribeDeleteBucket(t *testing.T) {
	var subscribeDeleteBucketTestData = []struct {
		ctx context.Context
		msg *deletebucket.DeleteBucketMsg
		out error
	}{
		{
			ctx,
			&deletebucket.DeleteBucketMsg{
				Index: "test",
			},
			nil,
		},
	}

	for _, tt := range subscribeDeleteBucketTestData {
		result := deleteBuckethandler.SubscribeDeleteBucket(tt.ctx, tt.msg)
		if result != tt.out {
			t.Errorf("Expected %v, got: %v", tt.out, result)
		}
	}
}
