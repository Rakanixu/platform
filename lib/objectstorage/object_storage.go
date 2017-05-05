package objectstorage

import (
	"io"
)

type ObjectStorage interface {
	Init() error
	CreateBucket(bucketName string) error
	DeleteBucket(bucketName string) error
	Upload(r io.ReadCloser, bucketName, key string) error
	Download(bucketName, key string, opts ...string) (io.ReadCloser, error)
	Delete(bucketName, key string) error
	SignedObjectStorageURL(bucketName string, objName string) (string, error)
}

var objStorage ObjectStorage

func Register(os ObjectStorage) {
	objStorage = os
}

func Init() error {
	return objStorage.Init()
}

func CreateBucket(bucketName string) error {
	return objStorage.CreateBucket(bucketName)
}

func DeleteBucket(bucketName string) error {
	return objStorage.DeleteBucket(bucketName)
}

func Upload(r io.ReadCloser, bucketName, key string) error {
	return objStorage.Upload(r, bucketName, key)
}

func Download(bucketName, key string, opts ...string) (io.ReadCloser, error) {
	return objStorage.Download(bucketName, key, opts...)
}

func Delete(bucketName, key string) error {
	return objStorage.Delete(bucketName, key)
}

func SignedObjectStorageURL(bucketName string, objName string) (string, error) {
	return objStorage.SignedObjectStorageURL(bucketName, objName)
}
