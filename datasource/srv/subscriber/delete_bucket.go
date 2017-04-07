package subscriber

import (
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	deletebucket "github.com/kazoup/platform/lib/protomsg/deletebucket"
	"golang.org/x/net/context"
)

func NewDeleteBucketHandler(cloudStorage *gcslib.GoogleCloudStorage) *deleteBucket {
	return &deleteBucket{
		cloudStorage: cloudStorage,
	}
}

type deleteBucket struct {
	cloudStorage *gcslib.GoogleCloudStorage
}

// SubscribeDeleteBucket subscribes to DeleteBucket Message to clean un a bicket in GC storage
func (d *deleteBucket) SubscribeDeleteBucket(ctx context.Context, msg *deletebucket.DeleteBucketMsg) error {
	return d.cloudStorage.DeleteBucket(msg.Index)
}
