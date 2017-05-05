package subscriber

import (
	"github.com/kazoup/platform/lib/objectstorage"
	deletebucket "github.com/kazoup/platform/lib/protomsg/deletebucket"
	"golang.org/x/net/context"
)

func NewDeleteBucketHandler() *deleteBucket {
	return new(deleteBucket)
}

type deleteBucket struct{}

// SubscribeDeleteBucket subscribes to DeleteBucket Message to clean un a bicket in GC storage
func (d *deleteBucket) SubscribeDeleteBucket(ctx context.Context, msg *deletebucket.DeleteBucketMsg) error {
	objectstorage.DeleteBucket(msg.Index)
	return objectstorage.DeleteBucket(msg.Index)
}
