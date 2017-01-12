package cloudstorage

// SignedObjectStorageURL
func (gcs *GoogleDriveCloudStorage) SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	return "", nil
}

// DeleteBucket deletes a bucket and all its contents in GCS account
func (gcs *GoogleDriveCloudStorage) DeleteBucket() error {
	return nil
}
