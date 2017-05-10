package objectstorage

import (
	"github.com/kazoup/platform/lib/globals"
	"github.com/kazoup/platform/lib/objectstorage"
	"io"
)

type MockStorage struct{}

func init() {
	objectstorage.Register(new(MockStorage))
}

func (ms *MockStorage) Init() error {
	return nil
}

func (ms *MockStorage) CreateBucket(bucketName string) error {
	return nil
}

func (ms *MockStorage) DeleteBucket(bucketName string) error {
	return nil
}

func (ms *MockStorage) Upload(r io.ReadCloser, bucketName, key string) error {
	return nil
}

func (ms *MockStorage) Download(bucketName, key string, opts ...string) (io.ReadCloser, error) {
	return readCloser{}, nil
}

func (ms *MockStorage) Delete(bucketName, key string) error {
	return nil
}

func (ms *MockStorage) SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	return globals.Mock, nil
}

// helper
type readCloser struct{}

func (rc readCloser) Read(p []byte) (int, error) {
	return 0, nil
}

func (rc readCloser) Close() (err error) {
	return nil
}
