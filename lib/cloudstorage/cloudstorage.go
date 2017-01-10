package cloudstorage

import "io"

// CloudStorage interface
type CloudStorage interface {
	CloudStorageOperations
	CloudStorageUtils
}

// CloudStorageOperations interface
type CloudStorageOperations interface {
	Upload(io.Reader, string) error
}

// CloudStorageUtils interface
type CloudStorageUtils interface {
	SignedObjectStorageURL(string) (string, error)
}
