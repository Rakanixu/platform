package cloudstorage

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"time"
)

// SignedObjectStorageURL
func (gcs *GoogleCloudStorage) SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	// Refresh when resource expires. 1 hour availability since was requested.
	SignedURLOption.Expires = time.Now().Add(3600000000000)

	url, err := SignedURL(bucketName, objName, SignedURLOption)
	if err != nil {
		return "", err
	}

	return url, nil
}

// DeleteBucket deletes a bucket and all its contents in GCS account
func (gcs *GoogleCloudStorage) DeleteBucket() error {
	c, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		return err
	}
	srv, err := storage.New(c)
	if err != nil {
		return err
	}

	_, err = srv.Buckets.Get(gcs.Endpoint.Index).Do()
	// Bucket does not exists
	if err != nil {
		return nil // If does not exists, it's already deleted, then success
	}

	r, err := srv.Objects.List(gcs.Endpoint.Index).Fields("items,nextPageToken").Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(gcs.Endpoint.Index, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gcs.deleteBucketNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	if err := srv.Buckets.Delete(gcs.Endpoint.Index).Do(); err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) deleteBucketNextPage(srv *storage.Service, nextPageToken string) error {
	r, err := srv.Objects.List(gcs.Endpoint.Index).Fields("items,nextPageToken").PageToken(nextPageToken).Do()
	if err != nil {
		return err
	}

	for _, v := range r.Items {
		if err := srv.Objects.Delete(gcs.Endpoint.Index, v.Name).Do(); err != nil {
			return err
		}
	}

	if len(r.NextPageToken) > 0 {
		if err := gcs.deleteBucketNextPage(srv, r.NextPageToken); err != nil {
			return err
		}
	}

	return nil
}
