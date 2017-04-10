package subscriber

import (
	gcslib "github.com/kazoup/platform/lib/googlecloudstorage"
	deletefilebucket "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	"golang.org/x/net/context"
)

func NewDeleteFileInBucketHandler(cloudStorage *gcslib.GoogleCloudStorage) *deleteFileInBucket {
	return &deleteFileInBucket{
		cloudStorage: cloudStorage,
	}
}

type deleteFileInBucket struct {
	cloudStorage *gcslib.GoogleCloudStorage
}

// SubscribeCleanBucket subscribes to DCleanBucket Message to remove thumbs not longer related with document in index
func (d *deleteFileInBucket) SubscribeDeleteFileInBucket(ctx context.Context, msg *deletefilebucket.DeleteFileInBucketMsg) error {
	return d.cloudStorage.Delete(msg.Index, msg.FileId)
}
