package subscriber

import (
	"github.com/kazoup/platform/lib/objectstorage"
	deletefilebucket "github.com/kazoup/platform/lib/protomsg/deletefileinbucket"
	"golang.org/x/net/context"
)

func NewDeleteFileInBucketHandler() *deleteFileInBucket {
	return new(deleteFileInBucket)
}

type deleteFileInBucket struct{}

// SubscribeCleanBucket subscribes to DCleanBucket Message to remove thumbs not longer related with document in index
func (d *deleteFileInBucket) SubscribeDeleteFileInBucket(ctx context.Context, msg *deletefilebucket.DeleteFileInBucketMsg) error {
	return objectstorage.Delete(msg.Index, msg.FileId)
}
